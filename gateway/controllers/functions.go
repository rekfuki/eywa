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

// GetFunctions returns list of functions scoped to the user
func GetFunctions(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)

	filter := k8s.LabelSelector().Equals(types.UserIDLabel, auth.UserID)
	fss, err := k8sClient.GetFunctionsStatusFiltered(filter)
	if err != nil {
		log.Errorf("Failed to get functions from k8s: ", err)
		return err
	}

	sfss := []types.FunctionStatusResponse{}
	for _, fs := range fss {
		sfss = append(sfss, makeFunctionStatusResponse(&fs))
	}

	return c.JSON(http.StatusOK, types.MultiFunctionStatusResponse{
		Objects: sfss,
		Total:   len(sfss),
	})
}

// GetFunction returns a specific service
func GetFunction(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)
	functionID := c.Param("function_id")

	filter := k8s.LabelSelector().
		Equals(types.FunctionIDLabel, functionID).
		Equals(types.UserIDLabel, auth.UserID)
	fs, err := k8sClient.GetFunctionStatusFiltered(filter)
	if err != nil {
		log.Errorf("Failed to get functions from k8s: ", err)
		return err
	}

	return c.JSON(http.StatusOK, makeFunctionStatusResponse(fs))
}

// DeployFunction deploys a new function onto k8s
func DeployFunction(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)
	rc := c.Get("registry").(*registry.Client)

	var dr types.DeployFunctionRequest
	if err := c.Bind(&dr); err != nil {
		return err
	}

	limits := k8sClient.GetLimits()
	errors := validateK8sParams(&dr.FunctionRequest, limits)
	if len(errors) > 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation error",
			"details": errors,
		})
	}

	serviceName := buildK8sName(dr.Name, auth.UserID)
	filter := k8s.LabelSelector().
		Equals(types.FunctionIDLabel, serviceName).
		Equals(types.UserIDLabel, auth.UserID)
	fs, err := k8sClient.GetFunctionStatusFiltered(filter)
	if err != nil {
		log.Errorf("Failed to retrieve function status: %s", err)
		return err
	}

	if fs != nil {
		return c.JSON(http.StatusBadRequest, "Function with specified name already exists")
	}

	filter = k8s.LabelSelector().
		In(types.SecretIDLabel, dr.Secrets).
		Equals(types.UserIDLabel, auth.UserID)
	secrets, err := k8sClient.GetSecretsFiltered(filter)
	if err != nil {
		log.Errorf("Failed to get secrets from k8s: %s", err)
		return err
	}

	notFoundSecrets := validateSecrets(dr.Secrets, secrets)
	if len(notFoundSecrets) > 0 {
		message := fmt.Sprintf("Following secrets not found: %#v", notFoundSecrets)
		return c.JSON(http.StatusNotFound, message)
	}

	if len(notFoundSecrets) > 0 {
		message := fmt.Sprintf("Following secrets not found: %#v", notFoundSecrets)
		return c.JSON(http.StatusNotFound, message)
	}

	image, err := rc.GetImage(dr.ImageID, auth.UserID)
	if err != nil {
		log.Errorf("Failed to get image from registry: %s", err)
		return err
	}

	if image == nil {
		return c.JSON(http.StatusNotFound, "Image Not Found")
	}

	dr.EnvVars = parseEnvVars(dr.FunctionRequest)

	fr := &k8s.DeployFunctionRequest{
		Image:         image.TaggedRegistry,
		Service:       serviceName,
		EnvVars:       dr.EnvVars,
		Secrets:       secrets,
		MinReplicas:   dr.MinReplicas,
		MaxReplicas:   dr.MaxReplicas,
		ScalingFactor: dr.ScalingFactor,
		Labels: map[string]string{
			types.UserIDLabel:          auth.UserID,
			types.ImageIDLabel:         image.ID,
			types.FunctionIDLabel:      serviceName,
			types.UserDefinedNameLabel: dr.Name,
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

	fs, err = k8sClient.DeployFunction(fr)
	if err != nil {
		log.Errorf("Failed to create function: %s", err)
		return err
	}

	return c.JSON(http.StatusCreated, makeFunctionStatusResponse(fs))
}

// UpdateFunction updates function deployment
func UpdateFunction(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)
	functionID := c.Param("function_id")

	var ur types.UpdateFunctionRequest
	if err := c.Bind(&ur); err != nil {
		return err
	}

	limits := k8sClient.GetLimits()
	errors := validateK8sParams(&ur.FunctionRequest, limits)
	if len(errors) > 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation error",
			"details": errors,
		})
	}

	filter := k8s.LabelSelector().
		Equals(types.FunctionIDLabel, functionID).
		Equals(types.UserIDLabel, auth.UserID)
	fs, err := k8sClient.GetFunctionStatusFiltered(filter)
	if err != nil {
		log.Errorf("Failed to retrieve function status: %s", err)
		return err
	}

	if fs == nil {
		return c.JSON(http.StatusBadRequest, "Function Not Found")
	}

	secrets := []k8s.Secret{}
	if len(ur.Secrets) > 0 {
		filter = k8s.LabelSelector().
			In(types.SecretIDLabel, ur.Secrets).
			Equals(types.UserIDLabel, auth.UserID)
		secrets, err = k8sClient.GetSecretsFiltered(filter)
		if err != nil {
			log.Errorf("Failed to get secrets from k8s: %s", err)
			return err
		}

		notFoundSecrets := validateSecrets(ur.Secrets, secrets)
		if len(notFoundSecrets) > 0 {
			message := fmt.Sprintf("Following secrets not found: %#v", notFoundSecrets)
			return c.JSON(http.StatusNotFound, message)
		}
	}

	ur.EnvVars = parseEnvVars(ur.FunctionRequest)

	fr := &k8s.DeployFunctionRequest{
		Image:         fs.Image,
		Service:       fs.Name,
		EnvVars:       ur.EnvVars,
		Secrets:       secrets,
		MinReplicas:   ur.MinReplicas,
		MaxReplicas:   ur.MaxReplicas,
		ScalingFactor: ur.ScalingFactor,
		Labels:        fs.Labels,
		Limits: &k8s.FunctionResources{
			CPU:    ur.Resources.MaxCPU,
			Memory: ur.Resources.MaxMemory,
		},
		Requests: &k8s.FunctionResources{
			CPU:    ur.Resources.MinCPU,
			Memory: ur.Resources.MinMemory,
		},
	}

	fs, err = k8sClient.UpdateFunction(fs.Name, fr)
	if err != nil {
		log.Errorf("Failed to update function: %s", err)
		return err
	}

	return c.JSON(http.StatusCreated, makeFunctionStatusResponse(fs))
}

// DeleteFunction deletes a function
func DeleteFunction(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)
	functionID := c.Param("function_id")

	filter := k8s.LabelSelector().
		Equals(types.FunctionIDLabel, functionID).
		Equals(types.UserIDLabel, auth.UserID)
	fs, err := k8sClient.GetFunctionStatusFiltered(filter)
	if err != nil {
		log.Errorf("Failed to to get function from k8s: %s", err)
		return err
	}

	if fs == nil {
		return c.JSON(http.StatusNotFound, "Function Not Found")
	} else if fs.DeletedAt != nil {
		return c.JSON(http.StatusBadRequest, "Function is terminating")
	}

	if err := k8sClient.DeleteFunction(fs.Name); err != nil {
		log.Errorf("Failed to delete function from k8s: %s", err)
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func buildK8sName(name, userID string) string {
	userUUID := uuid.FromStringOrNil(userID)
	functionID := uuid.NewV5(userUUID, name).String()
	return uuid.NewV5(userUUID, functionID).String()
}

func parseEnvVars(fr types.FunctionRequest) map[string]string {
	envVars := map[string]string{}
	fr.EnvVars["write_debug"] = "false"
	if fr.WriteDebug {
		envVars["write_debug"] = "true"
	}

	// Correct values should be validated by swagger
	rt, _ := time.ParseDuration(fr.ReadTimeout)
	if rt != time.Duration(0) {
		envVars["read_timeout"] = fr.ReadTimeout
	}

	wt, _ := time.ParseDuration(fr.WriteTimeout)
	if wt != time.Duration(0) {
		envVars["write_timeout"] = fr.WriteTimeout
	}

	envVars["max_inflight"] = fmt.Sprint(fr.MaxInflight)

	return envVars
}

func validateSecrets(uSecrets []string, k8sSecrets []k8s.Secret) []string {
	mappedSecrets := map[string]struct{}{}
	for _, secret := range k8sSecrets {
		if val, exists := secret.Labels[types.SecretIDLabel]; exists {
			mappedSecrets[val] = struct{}{}
		}
	}

	notFoundSecrets := []string{}
	for _, secret := range uSecrets {
		if _, exists := mappedSecrets[secret]; !exists {
			notFoundSecrets = append(notFoundSecrets, secret)
		}
	}

	return notFoundSecrets
}

// Validate further validates the payload after initial swagger validation
func validateK8sParams(dr *types.FunctionRequest, l *k8s.ResourceLimits) map[string][]string {
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

func makeFunctionStatusResponse(fs *k8s.FunctionStatus) (r types.FunctionStatusResponse) {
	r = types.FunctionStatusResponse{
		EnvVars:           fs.Env,
		Secrets:           fs.MountedSecrets,
		AvailableReplicas: fs.AvailableReplicas,
		MinReplicas:       fs.MinReplicas,
		MaxReplicas:       fs.MaxReplicas,
		ScalingFactor:     fs.ScalingFactor,
		CreatedAt:         fs.CreatedAt,
		UpdatedAt:         fs.UpdatedAt,
		DeletedAt:         fs.DeletedAt,
		Resources: types.FunctionResources{
			MaxCPU:    fs.Limits.CPU,
			MaxMemory: fs.Limits.Memory,
			MinCPU:    fs.Requests.CPU,
			MinMemory: fs.Requests.Memory,
		},
	}

	for k, v := range fs.Labels {
		switch k {
		case types.FunctionIDLabel:
			r.ID = v
		case types.ImageIDLabel:
			r.ImageID = v
		case types.UserDefinedNameLabel:
			r.Name = v
		}
	}

	for k, v := range fs.Env {
		switch k {
		case "max_inflight":
			i, err := strconv.Atoi(v)
			if err != nil {
				log.Errorf("Function %q has invalid max_inflight set %q: %s", fs.Name, v, err)
				continue
			}
			r.MaxInflight = i
		case "write_debug":
			if v == "true" {
				r.WriteDebug = true
			} else {
				r.WriteDebug = false
			}
		case "read_timeout":
			r.ReadTimeout = v
		case "write_timeout":
			r.WriteTimeout = v
		}
	}

	return
}
