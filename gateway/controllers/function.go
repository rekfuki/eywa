package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"eywa/gateway/clients/k8s"
	"eywa/gateway/clients/registry"
	"eywa/gateway/types"
	"eywa/go-libs/auth"
)

// GetServices returns list of services
func GetServices(c echo.Context) error {
	// TODO: When deployment is being deleted, it should show it as being terminated
	// Instead of non-existing because otherwise k8s will yell at us saying the deployment still exists
	return nil
}

// GetService returns a specific service
func GetService(c echo.Context) error {
	// TODO: When deployment is being deleted, it should show it as being terminated
	// Instead of non-existing because otherwise k8s will yell at us saying the deployment still exists
	return nil
}

// DeployFunction deploys a new function
func DeployFunction(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)
	rc := c.Get("registry").(*registry.Client)

	var dr types.DeployFunctionRequest
	if err := c.Bind(&dr); err != nil {
		return err
	}

	limits := k8sClient.GetLimits()
	errors := validateDeployRequest(&dr, limits)
	if len(errors) > 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation error",
			"details": errors,
		})
	}

	image, err := rc.GetImage(dr.ImageID, auth.UserID)
	if err != nil {
		log.Errorf("Failed to get image from registry: %s", err)
		return err
	}

	if image == nil {
		return c.JSON(http.StatusNotFound, "Image Not Found")
	}

	serviceName := fmt.Sprintf("%s-%s", dr.Name, image.ID)
	fn, err := k8sClient.GetFunctionStatus(serviceName)
	if err != nil {
		log.Errorf("Failed to retrieve function status: %s", err)
		return err
	}

	if fn != nil {
		return c.JSON(http.StatusBadRequest, "Function with specified name already exists")
	}

	dr.EnvVars["write_debug"] = "false"
	if dr.WriteDebug {
		dr.EnvVars["write_debug"] = "true"
	}

	if dr.ReadTimeout != time.Duration(0) {
		dr.EnvVars["read_timeout"] = dr.ReadTimeout.String()
	}

	if dr.WriteTimeout != time.Duration(0) {
		dr.EnvVars["write_timeout"] = dr.WriteTimeout.String()
	}

	dr.EnvVars["max_inflight"] = fmt.Sprint(dr.MaxInflight)

	unixStr := fmt.Sprint(time.Now().Unix())
	fr := &k8s.DeployFunctionRequest{
		Image:         image.TaggedRegistry,
		Service:       serviceName,
		EnvVars:       dr.EnvVars,
		Secrets:       dr.Secrets,
		MinReplicas:   dr.MinReplicas,
		MaxReplicas:   dr.MaxReplicas,
		ScalingFactor: dr.ScalingFactor,
		Labels: map[string]string{
			"user_id":           auth.UserID,
			"image_id":          image.ID,
			"function_id":       uuid.NewV4().String(),
			"user_defined_name": dr.Name,
			"created_at":        unixStr,
			"updated_at":        unixStr,
		},
		Limits: &k8s.FunctionResources{
			CPU:    dr.Resources.MaxCPU,
			Memory: dr.Resources.MaxMemory,
		},
		Requests: &k8s.FunctionResources{
			CPU:    dr.Resources.MinCPU,
			Memory: dr.Resources.MinMemory,
		},
	}

	fs, err := k8sClient.DeployFunction(fr, []k8s.Secret{})
	if err != nil {
		log.Errorf("Failed to deploy function: %s", err)
		return err
	}

	return c.JSON(http.StatusCreated, santiseFunctionStatus(fs))
}

// DeleteFunction deletes a function
func DeleteFunction(c echo.Context) error {
	return nil
}

// Validate further validates the payload after initial swagger validation
func validateDeployRequest(dr *types.DeployFunctionRequest, l *k8s.ResourceLimits) map[string][]string {
	errors := map[string][]string{}
	if dr.MaxReplicas < dr.MinReplicas {
		errors["max_replicas"] = append(errors["max_replicas"], "value must be at least equal to min_replicas")
	}

	if lt(dr.Resources.MaxCPU, dr.Resources.MinCPU) {
		errors["max_cpu"] = append(errors["max_cpu"], "value must be at least equal to min_cpu")
	}

	if lt(dr.Resources.MaxMemory, dr.Resources.MinMemory) {
		errors["max_memory"] = append(errors["max_memory"], "value must be at least equal to min_memory")
	}

	if lt(dr.Resources.MinCPU, l.MinCPU) {
		errors["min_cpu"] = append(errors["min_cpu"], "minimum value allowed "+l.MinCPU)
	}

	if gt(dr.Resources.MaxCPU, l.MaxCPU) {
		errors["max_cpu"] = append(errors["max_cpu"], "maximum value allowed "+l.MaxCPU)
	}

	if lt(dr.Resources.MinMemory, l.MinMem) {
		errors["min_memory"] = append(errors["min_memory"], "minimum value allowed "+l.MinMem)
	}

	if gt(dr.Resources.MaxMemory, l.MaxMem) {
		errors["max_memory"] = append(errors["max_memory"], "maximum value allowed "+l.MaxMem)
	}

	return errors
}

func gt(a, b string) bool {
	return cmpLimitStr(a, b)
}

func lt(a, b string) bool {
	return cmpLimitStr(b, a)
}

func cmpLimitStr(a, b string) bool {
	strValA := strings.Split(strings.ToLower(a), "m")[0]
	strValB := strings.Split(strings.ToLower(b), "m")[0]
	valA, _ := strconv.Atoi(strValA)
	valB, _ := strconv.Atoi(strValB)

	return valA > valB
}

func santiseFunctionStatus(fs *k8s.FunctionStatus) (r *types.FunctionStatusResponse) {
	r = &types.FunctionStatusResponse{
		FullName:          fs.Name,
		EnvVars:           fs.Env,
		MountedSecrets:    fs.MountedSecrets,
		AvailableReplicas: fs.AvailableReplicas,
		MinReplicas:       fs.MinReplicas,
		MaxReplicas:       fs.MaxReplicas,
		ScalingFactor:     fs.ScalingFactor,
		Resources: types.FunctionResources{
			MaxCPU:    fs.Limits.CPU,
			MaxMemory: fs.Limits.Memory,
			MinCPU:    fs.Requests.CPU,
			MinMemory: fs.Requests.Memory,
		},
	}

	for k, v := range fs.Labels {
		switch k {
		case "function_id":
			r.ID = v
		case "image_id":
			r.ImageID = v
		case "user_defined_name":
			r.ShortName = v
		case "created_at", "updated_at":
			t, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				log.Errorf("Function %q has invalid %q set %q: %s", r.FullName, k, v, err)
				continue
			}

			ut := time.Unix(t, 0)
			if k == "created_at" {
				r.CreatedAt = ut
			} else {
				r.UpdatedAt = ut
			}
		}
	}

	for k, v := range fs.Env {
		switch k {
		case "max_inflight":
			i, err := strconv.Atoi(v)
			if err != nil {
				log.Errorf("Function %q has invalid max_inflight set %q: %s", r.FullName, v, err)
				continue
			}
			r.MaxInflight = i
		case "write_debug":
			if v == "true" {
				r.WriteDebug = true
			} else {
				r.WriteDebug = false
			}
		case "read_timeout", "write_timeout":
			d, err := time.ParseDuration(v)
			if err != nil {
				log.Errorf("Function %q has invalid %q set %q: %s", r.FullName, k, v, err)
				continue
			}

			if k == "read_timeout" {
				r.ReadTimeout = d
			} else {
				r.WriteTimeout = d
			}
		}
	}

	return
}
