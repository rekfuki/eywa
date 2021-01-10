package controllers

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"eywa/go-libs/auth"
	"eywa/registry/builder"
	"eywa/registry/clients/docker"
	"eywa/registry/clients/mongo"
	"eywa/registry/types"
)

// GetImages returns all the images a user can access
func GetImages(c echo.Context) error {
	mc := c.Get("mongo").(*mongo.Client)
	auth := c.Get("auth").(*auth.Auth)
	page := c.Get("page_number").(int)
	perPage := c.Get("per_page").(int)

	total, images, err := mc.GetImages(auth.UserID, page, perPage)
	if err != nil {
		log.Errorf("Failed to retrieve images: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, types.GetImagesResponse{
		Objects: images,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	})
}

// GetImage returns a specific image
func GetImage(c echo.Context) error {
	mc := c.Get("mongo").(*mongo.Client)
	auth := c.Get("auth").(*auth.Auth)
	imageID := c.Param("image_id")

	image, err := mc.GetImage(imageID, auth.UserID)
	if err != nil {
		log.Errorf("Failed to retrieve image: %s", err)
		return err
	}

	if image == nil {
		return c.JSON(http.StatusNotFound, "Not Found")
	}

	if !auth.IsOperator() {
		image.TaggedRegistry = ""
	}

	return c.JSON(http.StatusOK, image)
}

// CreateImage handles image upserting
func CreateImage(c echo.Context) error {
	mc := c.Get("mongo").(*mongo.Client)
	builder := c.Get("builder").(*builder.Client)
	auth := c.Get("auth").(*auth.Auth)

	file, err := c.FormFile("source")
	if err != nil {
		log.Errorf("Failed to get source from payload: %s", err)
		return err
	}

	language := strings.ToLower(c.FormValue("language"))
	version := c.FormValue("version")
	name := c.FormValue("name")

	fullName := fmt.Sprintf("%s##%s##%s", language, name, version)
	id := uuid.NewV5(uuid.FromStringOrNil(auth.UserID), fullName).String()

	existingImage, err := mc.GetImage(id, auth.UserID)
	if err != nil {
		log.Errorf("Failed to retrieve image from db: %s", err)
		return err
	}

	if existingImage != nil {
		return c.JSON(http.StatusBadRequest, "Exact same image already exists")
	}

	src, err := file.Open()
	if err != nil {
		log.Errorf("Failed to open file header for reading: %s", err)
		return err
	}
	defer src.Close()

	body, err := ioutil.ReadAll(src)
	if err != nil {
		log.Errorf("Failed to read body: %s", err)
		return err
	}

	r, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Errorf("Failed to read zip file: %s", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	taggedRegistry := builder.Enqueue(types.BuildRequest{
		ID:        id,
		Language:  language,
		Version:   version,
		ZipReader: r,
	})

	image := types.Image{
		ID:             id,
		UserID:         auth.UserID,
		TaggedRegistry: taggedRegistry,
		Language:       language,
		Name:           name,
		Version:        version,
		CreatedAt:      time.Now(),
		State:          types.StateBuilding,
		Source:         base64.StdEncoding.EncodeToString(body),
	}

	if err := mc.CreateImage(image); err != nil {
		log.Errorf("Failed to create image in db: %s", err)
		return err
	}

	return c.JSON(http.StatusOK, image)
}

// DeleteImage deletes the image from db and registry
func DeleteImage(c echo.Context) error {
	mc := c.Get("mongo").(*mongo.Client)
	dc := c.Get("docker").(*docker.Client)
	auth := c.Get("auth").(*auth.Auth)
	imageID := c.Param("image_id")

	image, err := mc.GetImage(imageID, auth.UserID)
	if err != nil {
		log.Errorf("Failed to retrieve image: %s", err)
		return err
	}

	if image == nil {
		return c.JSON(http.StatusNotFound, "Not Found")
	}

	if err := dc.DeleteImage(image.ID, image.Version); err != nil {
		log.Errorf("Failed to delete image from docker: %s", err)
		return err
	}

	if err := mc.DeleteImage(imageID, auth.UserID); err != nil {
		log.Errorf("Failed to delete image form db: %s", err)
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
