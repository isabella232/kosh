package types

import "time"

// generated by "schematyper -o types/RequestType_BuildAddOrganization.go --package=types --ptr-for-omit BuildAddOrganization.json" -- DO NOT EDIT

type BuildAddOrganization struct {
	OrganizationID UUID `json:"organization_id"`
	Role           Role `json:"role"`
}

// generated by "schematyper -o types/RequestType_BuildAddUser.go --package=types --ptr-for-omit BuildAddUser.json" -- DO NOT EDIT

type BuildAddUser struct {
	Email  EmailAddress `json:"email,omitempty"`
	Role   Role         `json:"role"`
	UserID UUID         `json:"user_id,omitempty"`
}

// generated by "schematyper -o types/RequestType_BuildCreateDevices.go --package=types --ptr-for-omit BuildCreateDevices.json" -- DO NOT EDIT

type BuildCreateDevice struct {
	AssetTag     interface{}             `json:"asset_tag,omitempty"`
	ID           UUID                    `json:"id,omitempty"`
	Links        []Link                  `json:"links,omitempty"`
	SerialNumber DeviceSerialNumber      `json:"serial_number,omitempty"`
	Sku          MojoStandardPlaceholder `json:"sku"`
}

type BuildCreateDevices []BuildCreateDevice

// generated by "schematyper -o types/RequestType_BuildCreate.go --package=types --ptr-for-omit BuildCreate.json" -- DO NOT EDIT

type Admin struct {
	Email  EmailAddress `json:"email,omitempty"`
	UserID UUID         `json:"user_id,omitempty"`
}

type BuildCreate struct {
	Admins      []Admin                 `json:"admins,omitempty"`
	BuildID     UUID                    `json:"build_id,omitempty"`
	Description NonEmptyString          `json:"description,omitempty"`
	Name        MojoStandardPlaceholder `json:"name"`
	Started     time.Time               `json:"started,omitempty"`
}

// generated by "schematyper -o types/RequestType_BuildOrganizations.go --package=types --ptr-for-omit BuildOrganizations.json" -- DO NOT EDIT

type BuildOrganizations interface{}

// generated by "schematyper -o types/RequestType_BuildUpdate.go --package=types --ptr-for-omit BuildUpdate.json" -- DO NOT EDIT

type BuildUpdate struct {
	Completed   interface{}             `json:"completed,omitempty"`
	Description interface{}             `json:"description,omitempty"`
	Name        MojoStandardPlaceholder `json:"name,omitempty"`
	Started     interface{}             `json:"started,omitempty"`
}

// generated by "schematyper -o types/RequestType_DatacenterCreate.go --package=types --ptr-for-omit DatacenterCreate.json" -- DO NOT EDIT

type DatacenterCreate struct {
	Location   NonEmptyString `json:"location"`
	Region     NonEmptyString `json:"region"`
	Vendor     NonEmptyString `json:"vendor"`
	VendorName NonEmptyString `json:"vendor_name,omitempty"`
}

// generated by "schematyper -o types/RequestType_DatacenterUpdate.go --package=types --ptr-for-omit DatacenterUpdate.json" -- DO NOT EDIT

type DatacenterUpdate struct {
	Location   NonEmptyString `json:"location,omitempty"`
	Region     NonEmptyString `json:"region,omitempty"`
	Vendor     NonEmptyString `json:"vendor,omitempty"`
	VendorName NonEmptyString `json:"vendor_name,omitempty"`
}

// generated by "schematyper -o types/RequestType_DeviceLinks.go --package=types --ptr-for-omit DeviceLinks.json" -- DO NOT EDIT

type DeviceLinks struct {
	Links []Link `json:"links"`
}

// generated by "schematyper -o RequestType_DeviceReport.go --package=types --ptr-for-omit DeviceReport.json" -- DO NOT EDIT

type CpusItem map[string]interface{}

// the contents of a posted device report from relays and reporters
type DeviceReport struct {
	BiosVersion  string               `json:"bios_version"`
	Cpus         []CpusItem           `json:"cpus,omitempty"`
	DeviceType   string               `json:"device_type,omitempty"`
	Dimms        []Dimm               `json:"dimms,omitempty"`
	Disks        map[string]Disk      `json:"disks,omitempty"`
	Interfaces   map[string]Interface `json:"interfaces,omitempty"`
	Links        []Link               `json:"links,omitempty"`
	Os           *Os                  `json:"os,omitempty"`
	ProductName  string               `json:"product_name"`
	Relay        *Relay               `json:"relay,omitempty"`
	SerialNumber DeviceSerialNumber   `json:"serial_number"`
	Sku          string               `json:"sku"`
	SystemUUID   NonZeroUUID          `json:"system_uuid"`
	Temp         *Temp                `json:"temp,omitempty"`
	UptimeSince  string               `json:"uptime_since,omitempty"`
}

type Dimm struct {
	MemoryLocator      string      `json:"memory-locator"`
	MemorySerialNumber interface{} `json:"memory-serial-number,omitempty"`
	MemorySize         interface{} `json:"memory-size,omitempty"`
}

type Disk struct {
	BlockSz   int             `json:"block_sz,omitempty"`
	DriveType string          `json:"drive_type,omitempty"`
	Enclosure IntOrStringyInt `json:"enclosure,omitempty"`
	Firmware  string          `json:"firmware,omitempty"`
	Hba       IntOrStringyInt `json:"hba,omitempty"`
	Health    string          `json:"health,omitempty"`
	Model     string          `json:"model,omitempty"`
	Size      int             `json:"size,omitempty"`
	Slot      IntOrStringyInt `json:"slot,omitempty"`
	Temp      IntOrStringyInt `json:"temp,omitempty"`
	Transport string          `json:"transport,omitempty"`
	Vendor    string          `json:"vendor,omitempty"`
}

type DiskSerialNumber string

// an integer that may be presented as a json string
type IntOrStringyInt interface{}

type Interface struct {
	Ipaddr  interface{} `json:"ipaddr,omitempty"`
	Mac     Macaddr     `json:"mac"`
	Mtu     interface{} `json:"mtu,omitempty"`
	PeerMac interface{} `json:"peer_mac,omitempty"`
	Product string      `json:"product"`
	State   interface{} `json:"state,omitempty"`
	Vendor  string      `json:"vendor"`
}

type NonZeroUUID interface{}

type NonZeroUUIDEmbedded1 interface{}

type Os struct {
	Hostname string `json:"hostname"`
}

type Temp struct {
	Cpu0    IntOrStringyInt `json:"cpu0"`
	Cpu1    IntOrStringyInt `json:"cpu1"`
	Exhaust IntOrStringyInt `json:"exhaust,omitempty"`
	Inlet   IntOrStringyInt `json:"inlet,omitempty"`
}

// generated by "schematyper -o types/RequestType_HardwareProductCreate.go --package=types --ptr-for-omit HardwareProductCreate.json" -- DO NOT EDIT

type HardwareProductCreate struct {
	HardwareProductUpdate
	HardwareProductCreateEmbedded1
}

type HardwareProductCreateEmbedded1 interface{}

// generated by "schematyper -o types/RequestType_HardwareProductSpecification.go --package=types --ptr-for-omit HardwareProductSpecification.json" -- DO NOT EDIT

type Chassis struct {
	Memory *Memory `json:"memory,omitempty"`
}

type DiskSizeItem int

// this is the structure of the hardware_product.specification database column
type HardwareProductSpecification struct {
	Chassis  *Chassis                `json:"chassis,omitempty"`
	DiskSize map[string]DiskSizeItem `json:"disk_size,omitempty"`
}

type Memory struct {
	Dimms []Dimm `json:"dimms,omitempty"`
}

// generated by "schematyper -o types/RequestType_HardwareProductUpdate.go --package=types --ptr-for-omit HardwareProductUpdate.json" -- DO NOT EDIT

type HardwareProductUpdate struct {
	Alias             MojoStandardPlaceholder `json:"alias,omitempty"`
	BiosFirmware      string                  `json:"bios_firmware,omitempty"`
	CPUNum            int                     `json:"cpu_num,omitempty"`
	CPUType           string                  `json:"cpu_type,omitempty"`
	DimmsNum          int                     `json:"dimms_num,omitempty"`
	GenerationName    NonEmptyString          `json:"generation_name,omitempty"`
	HardwareVendorID  UUID                    `json:"hardware_vendor_id,omitempty"`
	HbaFirmware       interface{}             `json:"hba_firmware,omitempty"`
	LegacyProductName interface{}             `json:"legacy_product_name,omitempty"`
	Name              MojoStandardPlaceholder `json:"name,omitempty"`
	NicsNum           int                     `json:"nics_num,omitempty"`
	NvmeSsdNum        int                     `json:"nvme_ssd_num,omitempty"`
	NvmeSsdSize       interface{}             `json:"nvme_ssd_size,omitempty"`
	NvmeSsdSlots      interface{}             `json:"nvme_ssd_slots,omitempty"`
	Prefix            interface{}             `json:"prefix,omitempty"`
	PsuTotal          int                     `json:"psu_total,omitempty"`
	Purpose           string                  `json:"purpose,omitempty"`
	RAMTotal          int                     `json:"ram_total,omitempty"`
	RackUnitSize      PositiveInteger         `json:"rack_unit_size,omitempty"`
	RaidLunNum        int                     `json:"raid_lun_num,omitempty"`
	SasHddNum         int                     `json:"sas_hdd_num,omitempty"`
	SasHddSize        interface{}             `json:"sas_hdd_size,omitempty"`
	SasHddSlots       interface{}             `json:"sas_hdd_slots,omitempty"`
	SasSsdNum         int                     `json:"sas_ssd_num,omitempty"`
	SasSsdSize        interface{}             `json:"sas_ssd_size,omitempty"`
	SasSsdSlots       interface{}             `json:"sas_ssd_slots,omitempty"`
	SataHddNum        int                     `json:"sata_hdd_num,omitempty"`
	SataHddSize       interface{}             `json:"sata_hdd_size,omitempty"`
	SataHddSlots      interface{}             `json:"sata_hdd_slots,omitempty"`
	SataSsdNum        int                     `json:"sata_ssd_num,omitempty"`
	SataSsdSize       interface{}             `json:"sata_ssd_size,omitempty"`
	SataSsdSlots      interface{}             `json:"sata_ssd_slots,omitempty"`
	Sku               MojoStandardPlaceholder `json:"sku,omitempty"`
	Specification     interface{}             `json:"specification,omitempty"`
	UsbNum            int                     `json:"usb_num,omitempty"`
	ValidationPlanID  UUID                    `json:"validation_plan_id,omitempty"`
}

// generated by "schematyper -o types/RequestType_Login.go --package=types --ptr-for-omit Login.json" -- DO NOT EDIT

type Login struct {
	Email      EmailAddress   `json:"email,omitempty"`
	Password   NonEmptyString `json:"password"`
	SetSession bool           `json:"set_session,omitempty"`
	UserID     UUID           `json:"user_id,omitempty"`
}

// generated by "schematyper -o types/RequestType_NewUser.go --package=types --ptr-for-omit NewUser.json" -- DO NOT EDIT

type NewUser struct {
	Email    EmailAddress   `json:"email"`
	IsAdmin  bool           `json:"is_admin,omitempty"`
	Name     NonEmptyString `json:"name"`
	Password NonEmptyString `json:"password,omitempty"`
}

// generated by "schematyper -o types/RequestType_NewUserToken.go --package=types --ptr-for-omit NewUserToken.json" -- DO NOT EDIT

type NewUserToken struct {
	Name string `json:"name"`
}

// generated by "schematyper -o types/RequestType_OrganizationAddUser.go --package=types --ptr-for-omit OrganizationAddUser.json" -- DO NOT EDIT

type OrganizationAddUser struct {
	Email  EmailAddress `json:"email,omitempty"`
	Role   Role         `json:"role"`
	UserID UUID         `json:"user_id,omitempty"`
}

// generated by "schematyper -o types/RequestType_OrganizationCreate.go --package=types --ptr-for-omit OrganizationCreate.json" -- DO NOT EDIT

type OrganizationCreate struct {
	Admins      []Admin                 `json:"admins"`
	Description NonEmptyString          `json:"description,omitempty"`
	Name        MojoStandardPlaceholder `json:"name"`
}

// generated by "schematyper -o types/RequestType_OrganizationUpdate.go --package=types --ptr-for-omit OrganizationUpdate.json" -- DO NOT EDIT

type OrganizationUpdate struct {
	Description interface{}             `json:"description,omitempty"`
	Name        MojoStandardPlaceholder `json:"name,omitempty"`
}

// generated by "schematyper -o types/RequestType_RegisterRelay.go --package=types --ptr-for-omit RegisterRelay.json" -- DO NOT EDIT

type RegisterRelay struct {
	Ipaddr  string             `json:"ipaddr,omitempty"`
	Name    NonEmptyString     `json:"name,omitempty"`
	SSHPort NonNegativeInteger `json:"ssh_port,omitempty"`
	Serial  RelaySerialNumber  `json:"serial"`
	Version string             `json:"version,omitempty"`
}

// generated by "schematyper -o types/RequestType_UpdateUser.go --package=types --ptr-for-omit UpdateUser.json" -- DO NOT EDIT

type UpdateUser struct {
	Email   EmailAddress   `json:"email,omitempty"`
	IsAdmin bool           `json:"is_admin,omitempty"`
	Name    NonEmptyString `json:"name,omitempty"`
}
