// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package googlecompute

import (
	"fmt"

	compute "google.golang.org/api/compute/v1"
	oslogin "google.golang.org/api/oslogin/v1"
)

// DriverMock is a Driver implementation that is a mocked out so that
// it can be used for tests.
type DriverMock struct {
	CreateDiskConfig   BlockDevice
	CreateDiskResultCh <-chan *compute.Disk
	CreateDiskErrCh    <-chan error

	CreateImageProjectId        string
	CreateImageName             string
	CreateImageDesc             string
	CreateImageFamily           string
	CreateImageEncryptionKey    *compute.CustomerEncryptionKey
	CreateImageLabels           map[string]string
	CreateImageLicenses         []string
	CreateImageFeatures         []string
	CreateImageStorageLocations []string
	CreateImageZone             string
	CreateImageDisk             string
	CreateImageResultProjectId  string
	CreateImageResultSelfLink   string
	CreateImageResultSizeGb     int64
	CreateImageErrCh            <-chan error
	CreateImageResultCh         <-chan *Image

	DeleteProjectId  string
	DeleteImageName  string
	DeleteImageErrCh <-chan error

	DeleteInstanceZone  string
	DeleteInstanceName  string
	DeleteInstanceErrCh <-chan error
	DeleteInstanceErr   error

	DeleteDiskZone  string
	DeleteDiskName  string
	DeleteDiskErrCh chan error
	DeleteDiskErr   error

	GetDiskName   string
	GetDiskZone   string
	GetDiskResult *compute.Disk
	GetDiskErr    error

	GetImageName           string
	GetImageSourceProjects []string
	GetImageFromFamily     bool
	GetImageResult         *Image
	GetImageErr            error

	GetImageFromProjectProject    string
	GetImageFromProjectName       string
	GetImageFromProjectFromFamily bool
	GetImageFromProjectResult     *Image
	GetImageFromProjectErr        error

	GetInstanceMetadataZone   string
	GetInstanceMetadataName   string
	GetInstanceMetadataKey    string
	GetInstanceMetadataResult string
	GetInstanceMetadataErr    error

	GetNatIPZone   string
	GetNatIPName   string
	GetNatIPResult string
	GetNatIPErr    error

	GetInternalIPZone   string
	GetInternalIPName   string
	GetInternalIPResult string
	GetInternalIPErr    error

	GetSerialPortOutputZone   string
	GetSerialPortOutputName   string
	GetSerialPortOutputResult string
	GetSerialPortOutputErr    error

	ImageExistsProjectId string
	ImageExistsName      string
	ImageExistsResult    bool

	RunInstanceConfig *InstanceConfig
	RunInstanceErrCh  <-chan error
	RunInstanceErr    error

	CreateOrResetWindowsPasswordZone     string
	CreateOrResetWindowsPasswordInstance string
	CreateOrResetWindowsPasswordConfig   *WindowsPasswordConfig
	CreateOrResetWindowsPasswordErr      error
	CreateOrResetWindowsPasswordErrCh    <-chan error

	WaitForInstanceState string
	WaitForInstanceZone  string
	WaitForInstanceName  string
	WaitForInstanceErrCh <-chan error

	AddToInstanceMetadataZone    string
	AddToInstanceMetadataName    string
	AddToInstanceMetadataKVPairs map[string]string
	AddToInstanceMetadataErrCh   <-chan error
	AddToInstanceMetadataErr     error
}

func (d *DriverMock) CreateImage(project, name, description, family, zone, disk string, image_labels map[string]string, image_licenses []string, image_features []string, image_encryption_key *compute.CustomerEncryptionKey, imageStorageLocations []string) (<-chan *Image, <-chan error) {
	d.CreateImageProjectId = project
	d.CreateImageName = name
	d.CreateImageDesc = description
	d.CreateImageFamily = family
	d.CreateImageLabels = image_labels
	d.CreateImageLicenses = image_licenses
	d.CreateImageFeatures = image_features
	d.CreateImageStorageLocations = imageStorageLocations
	d.CreateImageZone = zone
	d.CreateImageDisk = disk
	d.CreateImageEncryptionKey = image_encryption_key
	if d.CreateImageResultProjectId == "" {
		d.CreateImageResultProjectId = "test"
	}
	if d.CreateImageResultSelfLink == "" {
		d.CreateImageResultSelfLink = fmt.Sprintf(
			"http://content.googleapis.com/compute/v1/%s/global/licenses/test",
			d.CreateImageResultProjectId)
	}
	if d.CreateImageResultSizeGb == 0 {
		d.CreateImageResultSizeGb = 10
	}
	imageFeatures := make([]*compute.GuestOsFeature, 0, len(image_features))
	for _, v := range image_features {
		imageFeatures = append(imageFeatures, &compute.GuestOsFeature{
			Type: v,
		})
	}
	resultCh := d.CreateImageResultCh
	if resultCh == nil {
		ch := make(chan *Image, 1)
		ch <- &Image{
			GuestOsFeatures: imageFeatures,
			Labels:          d.CreateImageLabels,
			Licenses:        d.CreateImageLicenses,
			Name:            name,
			ProjectId:       d.CreateImageResultProjectId,
			SelfLink:        d.CreateImageResultSelfLink,
			SizeGb:          d.CreateImageResultSizeGb,
		}
		close(ch)
		resultCh = ch
	}

	errCh := d.CreateImageErrCh
	if errCh == nil {
		ch := make(chan error)
		close(ch)
		errCh = ch
	}

	return resultCh, errCh
}

func (d *DriverMock) DeleteImage(project, name string) <-chan error {
	d.DeleteProjectId = project
	d.DeleteImageName = name

	resultCh := d.DeleteImageErrCh
	if resultCh == nil {
		ch := make(chan error)
		close(ch)
		resultCh = ch
	}

	return resultCh
}

func (d *DriverMock) DeleteInstance(zone, name string) (<-chan error, error) {
	d.DeleteInstanceZone = zone
	d.DeleteInstanceName = name

	resultCh := d.DeleteInstanceErrCh
	if resultCh == nil {
		ch := make(chan error)
		close(ch)
		resultCh = ch
	}

	return resultCh, d.DeleteInstanceErr
}

func (d *DriverMock) CreateDisk(diskConfig BlockDevice) (<-chan *compute.Disk, <-chan error) {
	d.CreateDiskConfig = diskConfig

	resultCh := d.CreateDiskResultCh
	if resultCh == nil {
		ch := make(chan *compute.Disk)
		close(ch)
		resultCh = ch
	}

	errCh := d.CreateDiskErrCh
	if errCh != nil {
		ch := make(chan error)
		close(ch)
		errCh = ch
	}

	return resultCh, errCh
}

func (d *DriverMock) DeleteDisk(zone, name string) <-chan error {
	d.DeleteDiskZone = zone
	d.DeleteDiskName = name

	resultCh := d.DeleteDiskErrCh
	if resultCh == nil {
		ch := make(chan error)
		resultCh = ch
	}

	if d.DeleteDiskErr != nil {
		resultCh <- d.DeleteDiskErr
	}

	close(resultCh)

	return resultCh
}

func (d *DriverMock) GetDisk(zoneOrRegion, name string) (*compute.Disk, error) {
	d.GetDiskZone = zoneOrRegion
	d.GetDiskName = name

	return d.GetDiskResult, d.GetDiskErr
}

func (d *DriverMock) GetImage(name string, fromFamily bool) (*Image, error) {
	d.GetImageName = name
	d.GetImageFromFamily = fromFamily
	return d.GetImageResult, d.GetImageErr
}
func (d *DriverMock) GetImageFromProjects(projects []string, name string, fromFamily bool) (*Image, error) {
	d.GetImageSourceProjects = projects
	d.GetImageFromProjectName = name
	d.GetImageFromProjectFromFamily = fromFamily
	return d.GetImageFromProjectResult, d.GetImageFromProjectErr
}

func (d *DriverMock) GetImageFromProject(project, name string, fromFamily bool) (*Image, error) {
	d.GetImageFromProjectProject = project
	d.GetImageFromProjectName = name
	d.GetImageFromProjectFromFamily = fromFamily
	return d.GetImageFromProjectResult, d.GetImageFromProjectErr
}

func (d *DriverMock) GetInstanceMetadata(zone, name, key string) (string, error) {
	d.GetInstanceMetadataZone = zone
	d.GetInstanceMetadataName = name
	d.GetInstanceMetadataKey = key
	return d.GetInstanceMetadataResult, d.GetInstanceMetadataErr
}

func (d *DriverMock) GetNatIP(zone, name string) (string, error) {
	d.GetNatIPZone = zone
	d.GetNatIPName = name
	return d.GetNatIPResult, d.GetNatIPErr
}

func (d *DriverMock) GetInternalIP(zone, name string) (string, error) {
	d.GetInternalIPZone = zone
	d.GetInternalIPName = name
	return d.GetInternalIPResult, d.GetInternalIPErr
}

func (d *DriverMock) GetSerialPortOutput(zone, name string) (string, error) {
	d.GetSerialPortOutputZone = zone
	d.GetSerialPortOutputName = name
	return d.GetSerialPortOutputResult, d.GetSerialPortOutputErr
}

func (d *DriverMock) ImageExists(project, name string) bool {
	d.ImageExistsProjectId = project
	d.ImageExistsName = name
	return d.ImageExistsResult
}

func (d *DriverMock) RunInstance(c *InstanceConfig) (<-chan error, error) {
	d.RunInstanceConfig = c

	resultCh := d.RunInstanceErrCh
	if resultCh == nil {
		ch := make(chan error)
		close(ch)
		resultCh = ch
	}

	return resultCh, d.RunInstanceErr
}

func (d *DriverMock) WaitForInstance(state, zone, name string) <-chan error {
	d.WaitForInstanceState = state
	d.WaitForInstanceZone = zone
	d.WaitForInstanceName = name

	resultCh := d.WaitForInstanceErrCh
	if resultCh == nil {
		ch := make(chan error)
		close(ch)
		resultCh = ch
	}

	return resultCh
}

func (d *DriverMock) GetWindowsPassword() (string, error) {
	return "", nil
}

func (d *DriverMock) CreateOrResetWindowsPassword(instance, zone string, c *WindowsPasswordConfig) (<-chan error, error) {

	d.CreateOrResetWindowsPasswordInstance = instance
	d.CreateOrResetWindowsPasswordZone = zone
	d.CreateOrResetWindowsPasswordConfig = c

	c.password = "MOCK_PASSWORD"

	resultCh := d.CreateOrResetWindowsPasswordErrCh
	if resultCh == nil {
		ch := make(chan error)
		close(ch)
		resultCh = ch
	}

	return resultCh, d.CreateOrResetWindowsPasswordErr
}

func (d *DriverMock) ImportOSLoginSSHKey(user, key string) (*oslogin.LoginProfile, error) {
	account := oslogin.PosixAccount{Primary: true, Username: "testing_packer_io"}
	profile := oslogin.LoginProfile{
		PosixAccounts: []*oslogin.PosixAccount{&account},
	}
	return &profile, nil
}

func (d *DriverMock) DeleteOSLoginSSHKey(user, fingerprint string) error {
	return nil
}

func (d *DriverMock) AddToInstanceMetadata(zone string, name string, metadata map[string]string) error {
	d.AddToInstanceMetadataZone = zone
	d.AddToInstanceMetadataName = name
	d.AddToInstanceMetadataKVPairs = metadata

	resultCh := d.AddToInstanceMetadataErrCh
	if resultCh == nil {
		ch := make(chan error)
		close(ch)
	}

	return nil
}
