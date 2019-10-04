// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

//lint:file-ignore U1000 WIP

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	cli "github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
)

type Racks struct {
	*Conch
}

func (c *Conch) Racks() *Racks {
	return &Racks{c}
}

type RackList []Rack
type Rack struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	RoomID       uuid.UUID `json:"datacenter_room_id"`
	RoleID       uuid.UUID `json:"rack_role_id"`
	SerialNumber string    `json:"serial_number,omitempty"`
	AssetTag     string    `json:"asset_tag,omitempty"`
	Phase        string    `json:"phase"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	BuildID      uuid.UUID `json:"build_id"` // TODO build support

	Role RackRole `json:"-"`
	Room Room     `json:"-"`
}

func (rl RackList) String() string {
	if API.JsonOnly {
		return API.AsJSON(rl)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"ID",
		"Name",
		"Room",
		"Role",
		"Serial Number",
		"Asset Tag",
		"Phase",
		"Created",
		"Updated",
	})

	for _, r := range rl {
		var role string
		if (r.RoleID != uuid.UUID{}) {
			role = r.Role.Name
		}

		var room string
		if (r.RoomID != uuid.UUID{}) {
			room = r.Room.Alias
		}

		table.Append([]string{
			CutUUID(r.ID.String()),
			r.Name,
			room,
			role,
			r.SerialNumber,
			r.AssetTag,
			r.Phase,
			TimeStr(r.Created),
			TimeStr(r.Updated),
		})
	}

	table.Render()
	return tableString.String()

}

func (r Rack) String() string {
	if API.JsonOnly {
		return API.AsJSON(r)
	}

	t, err := NewTemplate().Parse(rackTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, r); err != nil {
		panic(err)
	}

	return buf.String()
}

func (r *Racks) GetAll() RackList {
	rl := make(RackList, 0)

	res := r.Do(r.Sling().Get("/rack"))
	if ok := res.Parse(&rl); !ok {
		panic(fmt.Sprintf("%v", res))
	}

	roles := make(map[uuid.UUID]RackRole)
	rooms := make(map[uuid.UUID]Room)

	list := make(RackList, 0)

	for _, rack := range rl {
		if (rack.RoleID != uuid.UUID{}) {
			if role, ok := roles[rack.RoleID]; ok {
				rack.Role = role
			} else {
				rack.Role = API.RackRoles().Get(rack.RoleID)
				roles[rack.RoleID] = rack.Role
			}
		}

		if (rack.RoomID != uuid.UUID{}) {
			if room, ok := rooms[rack.RoomID]; ok {
				rack.Room = room
			} else {
				rack.Room = API.Rooms().Get(rack.RoomID)
				rooms[rack.RoomID] = rack.Room
			}
		}

		list = append(list, rack)
	}

	return list
}

func (r *Racks) FindID(id string) (bool, uuid.UUID) {
	ids := make([]uuid.UUID, 0)
	for _, rack := range r.GetAll() {
		ids = append(ids, rack.ID)
	}

	return FindUUID(id, ids)
}

func (r *Racks) Get(id uuid.UUID) Rack {
	var rack Rack
	uri := fmt.Sprintf(
		"/rack/%s",
		url.PathEscape(id.String()),
	)

	res := r.Do(r.Sling().Get(uri))
	if ok := res.Parse(&rack); !ok {
		panic(res)
	}

	if (rack.RoleID != uuid.UUID{}) {
		rack.Role = API.RackRoles().Get(rack.RoleID)
	}

	if (rack.RoomID != uuid.UUID{}) {
		rack.Room = API.Rooms().Get(rack.RoomID)
	}

	return rack
}

func (r *Racks) Create(name string, roomID uuid.UUID, roleID uuid.UUID, phase string) Rack {
	payload := make(map[string]string)
	if name == "" {
		panic(errors.New("'name' cannot be empty"))
	}
	payload["name"] = name

	if (roomID == uuid.UUID{}) {
		panic(errors.New("'roomID' cannot be empty"))
	}
	payload["datacenter_room_id"] = roomID.String()

	if (roleID == uuid.UUID{}) {
		panic(errors.New("'roleID' cannot be empty"))
	}
	payload["rack_role_id"] = roleID.String()

	if phase != "" {
		payload["phase"] = phase
	}

	var rack Rack

	// We get a 303 on success
	res := r.Do(
		r.Sling().New().Post("/rack").
			Set("Content-Type", "application/json").
			BodyJSON(payload),
	)

	if ok := res.Parse(&rack); !ok {
		panic(res)
	}

	if (rack.RoleID != uuid.UUID{}) {
		rack.Role = API.RackRoles().Get(rack.RoleID)
	}

	return rack
}

func (r *Racks) Update(
	id uuid.UUID,
	newName string,
	roomID uuid.UUID,
	roleID uuid.UUID,
	phase string,
	serialNumber *string,
	assetTag *string,
) Rack {

	payload := make(map[string]interface{})
	if newName != "" {
		payload["name"] = newName
	}

	if (roomID != uuid.UUID{}) {
		payload["datacenter_room_id"] = roomID.String()
	}

	if (roleID != uuid.UUID{}) {
		payload["rack_role_id"] = roleID.String()
	}

	if phase != "" {
		payload["phase"] = phase
	}

	if serialNumber == nil {
		payload["serial_number"] = nil
	} else if *serialNumber != "" {
		payload["serial_number"] = *serialNumber
	}

	if assetTag == nil {
		payload["asset_tag"] = nil
	} else if *assetTag != "" {
		payload["asset_tag"] = *assetTag
	}

	if len(payload) == 0 {
		return r.Get(id)
	}

	var rack Rack

	uri := fmt.Sprintf(
		"/rack/%s",
		url.PathEscape(id.String()),
	)

	// We get a 303 on success
	res := r.Do(
		r.Sling().New().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(payload),
	)

	if ok := res.Parse(&rack); !ok {
		panic(res)
	}

	if (rack.RoleID != uuid.UUID{}) {
		rack.Role = API.RackRoles().Get(rack.RoleID)
	}

	return rack
}

func (r *Racks) Delete(id uuid.UUID) {
	uri := fmt.Sprintf(
		"/rack/%s",
		url.PathEscape(id.String()),
	)

	res := r.Do(r.Sling().New().Delete(uri))

	if res.StatusCode() != 204 {
		// I know this is weird. Like in other places, it should be impossible
		// to reach here unless the status code is 204. The API returns 204
		// (which gets us here) or 409 (which will explode before it gets here).
		// If we got here via some other code, then there's some new behavior
		// that we need to know about.

		panic(res)
	}
}

/****/

type RackLayout []RackLayoutSlot

func (r RackLayout) Len() int {
	return len(r)
}

func (r RackLayout) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RackLayout) Less(i, j int) bool {
	return r[i].RackUnitStart > r[j].RackUnitStart
}

func (rl RackLayout) String() string {
	sort.Sort(rl)
	if API.JsonOnly {
		return API.AsJSON(rl)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"Rack Unit Start",
		"Rack Unit Size",
		"ID",
		"Hardware Product",
		"Created",
		"Updated",
	})

	products := make(map[uuid.UUID]HardwareProduct)

	for _, r := range rl {
		var hpName = ""
		if (r.HardwareProductID != uuid.UUID{}) {
			var hp HardwareProduct
			if _, ok := products[r.HardwareProductID]; ok {
				hp = products[r.HardwareProductID]
			} else {
				hp = API.Hardware().GetProduct(r.HardwareProductID)
				products[r.HardwareProductID] = hp
			}

			hpName = fmt.Sprintf(
				"%s (%s)",
				hp.Alias,
				hp.Name,
			)

		}
		table.Append([]string{
			strconv.Itoa(r.RackUnitStart),
			strconv.Itoa(r.RackUnitSize),
			CutUUID(r.ID.String()),
			hpName,
			TimeStr(r.Created),
			TimeStr(r.Updated),
		})
	}

	table.Render()
	return tableString.String()
}

func (rl RackLayout) Export() string {
	type Slot struct {
		RU           int       `json:"ru_start"`
		ProductID    uuid.UUID `json:"product_id,omitempty"`
		ProductName  string    `json:"product_name,omitempty"`
		ProductAlias string    `json:"product_alias,omitempty"`
	}
	slots := make([]Slot, 0)

	sort.Sort(rl)

	hpCache := make(map[uuid.UUID]HardwareProduct)

	for _, slot := range rl {
		if _, ok := hpCache[slot.HardwareProductID]; !ok {
			hpCache[slot.HardwareProductID] = API.Hardware().GetProduct(slot.HardwareProductID)
		}

		slots = append(slots, Slot{
			slot.RackUnitStart,
			slot.HardwareProductID,
			hpCache[slot.HardwareProductID].Name,
			hpCache[slot.HardwareProductID].Alias,
		})
	}

	return API.AsJSON(slots)
}

func (r *Racks) ImportLayout(rackID uuid.UUID, b []byte) RackLayout {
	type Slot struct {
		RU           int       `json:"ru_start"`
		ProductID    uuid.UUID `json:"product_id,omitempty"`
		ProductName  string    `json:"product_name,omitempty"`
		ProductAlias string    `json:"product_alias,omitempty"`
	}

	imported := make([]Slot, 0)
	if err := json.Unmarshal(b, &imported); err != nil {
		panic(err)
	}

	hpCache := make(map[string]HardwareProduct)
	slots := make([]Slot, 0)

	for _, row := range imported {
		var slot Slot

		slot.RU = row.RU
		slot.ProductID = row.ProductID
		if (row.ProductID != uuid.UUID{}) {
			slots = append(slots, slot)
			continue
		}

		if row.ProductName != "" {
			if hp, ok := hpCache[row.ProductName]; ok {
				hpCache[row.ProductName] = hp
			} else {
				hpCache[row.ProductName] = API.Hardware().GetProductByName(row.ProductName)
			}
			slot.ProductID = hpCache[row.ProductName].ID
		} else if row.ProductAlias != "" {
			if hp, ok := hpCache[row.ProductAlias]; ok {
				hpCache[row.ProductAlias] = hp
			} else {
				hpCache[row.ProductAlias] = API.Hardware().GetProductByAlias(row.ProductAlias)
			}
			slot.ProductID = hpCache[row.ProductAlias].ID
		} else {
			panic(fmt.Errorf("RU %d entry does not have a product id, name, or alias", row.RU))
		}
		slots = append(slots, slot)
	}

	// There is no way to do this atomically. The api has no way to perform
	// this action other than deleting each row at a time and then putting them
	// back. If this seems really risky to you, then we are of the same mind.
	for _, row := range r.Layouts(rackID) {
		r.DeleteLayoutSlot(row.ID)
	}
	for _, slot := range slots {
		r.SaveLayoutSlot(rackID, slot.RU, slot.ProductID)
	}

	return r.Layouts(rackID)
}

type RackLayoutSlot struct {
	ID                uuid.UUID `json:"id"`
	RackID            uuid.UUID `json:"rack_id"`
	HardwareProductID uuid.UUID `json:"hardware_product_id"`
	RackUnitStart     int       `json:"rack_unit_start"`
	RackUnitSize      int       `json:"rack_unit_size"`
	Created           time.Time `json:"created"`
	Updated           time.Time `json:"updated"`
}

func (r *Racks) Layouts(id uuid.UUID) RackLayout {
	uri := fmt.Sprintf(
		"/rack/%s/layouts",
		url.PathEscape(id.String()),
	)

	layouts := make(RackLayout, 0)

	res := r.Do(r.Sling().New().Get(uri))
	if ok := res.Parse(&layouts); !ok {
		panic(res)
	}

	return layouts
}

func (r *Racks) DeleteLayoutSlot(id uuid.UUID) {
	uri := fmt.Sprintf(
		"/layout/%s",
		url.PathEscape(id.String()),
	)

	if res := r.Do(r.Sling().New().Delete(uri)); res.StatusCode() != 204 {
		panic(res)
	}
}

func (r *Racks) SaveLayoutSlot(rackID uuid.UUID, ruStart int, hardwareProductID uuid.UUID) (l RackLayoutSlot) {
	payload := make(map[string]interface{})
	payload["rack_id"] = rackID.String()
	payload["hardware_product_id"] = hardwareProductID.String()
	payload["rack_unit_start"] = ruStart

	res := r.Do(
		r.Sling().New().Post("/layout").
			Set("Content-Type", "application/json").
			BodyJSON(payload),
	)
	if ok := res.Parse(&l); !ok {
		panic(res)
	}
	return l
}

/****/

type RackAssignments []RackAssignment

func (r RackAssignments) Len() int {
	return len(r)
}

func (r RackAssignments) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RackAssignments) Less(i, j int) bool {
	return r[i].RackUnitStart > r[j].RackUnitStart
}

type RackAssignment struct {
	DeviceID            uuid.UUID `json:"device_id"`
	DeviceAssetTag      string    `json:"device_asset_tag,omitempty"`
	HardwareProductName string    `json:"hardware_product_name,omitempty"`
	RackUnitStart       int       `json:"rack_unit_start"`
	RackUnitSize        int       `json:"rack_unit_size"`
}

func (a RackAssignments) String() string {
	sort.Sort(a)
	if API.JsonOnly {
		return API.AsJSON(a)
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	TableToMarkdown(table)

	table.SetHeader([]string{
		"Device Serial",
		"Device Asset Tag",
		"Hardware Product",
		"Rack Unit Start",
		"Rack Unit Size",
	})

	for _, r := range a {
		var serial string
		if (r.DeviceID != uuid.UUID{}) {
			serial = API.Devices().Get(r.DeviceID.String()).Serial
		}

		table.Append([]string{
			serial,
			r.DeviceAssetTag,
			r.HardwareProductName,
			strconv.Itoa(r.RackUnitStart),
			strconv.Itoa(r.RackUnitSize),
		})
	}

	table.Render()
	return tableString.String()

}

func (r *Racks) Assignments(id uuid.UUID) RackAssignments {
	uri := fmt.Sprintf(
		"/rack/%s/assignment",
		url.PathEscape(id.String()),
	)

	assignments := make(RackAssignments, 0)
	res := r.Do(r.Sling().New().Get(uri))
	if ok := res.Parse(&assignments); !ok {
		panic(res)
	}
	return assignments
}

func (r *Racks) ImportAssignments(id uuid.UUID, b []byte) RackAssignments {
	type Assignment struct {
		DeviceID       uuid.UUID `json:"device_id"`
		RackUnitStart  int       `json:"rack_unit_start"`
		DeviceAssetTag string    `json:"device_asset_tag,omitempty"`
	}

	imported := make([]Assignment, 0)
	if err := json.Unmarshal(b, &imported); err != nil {
		panic(err)
	}

	uri := fmt.Sprintf(
		"/rack/%s/assignment",
		url.PathEscape(id.String()),
	)

	r.Do(
		r.Sling().New().Post(uri).
			Set("Content-Type", "application/json").
			BodyJSON(imported),
	)

	return r.Assignments(id)
}

/****/

func init() {
	App.Command("racks", "Work with datacenter racks", func(cmd *cli.Cmd) {
		cmd.Before = RequireSysAdmin
		cmd.Command("get", "Get a list of all racks", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Racks().GetAll())
			}
		})

		cmd.Command("create", "Create a new rack", func(cmd *cli.Cmd) {
			var (
				nameOpt      = cmd.StringOpt("name", "", "Name of the rack")
				roomAliasOpt = cmd.StringOpt("room", "", "Alias of the datacenter room")
				roleNameOpt  = cmd.StringOpt("role", "", "Name of the role")
				phaseOpt     = cmd.StringOpt("phase", "", "Optional phase for the rack")
			)

			cmd.Spec = "--name --room --role [OPTIONS]"
			cmd.Action = func() {
				var (
					roomID uuid.UUID
					roleID uuid.UUID
					ok     bool
				)

				// The user can be very silly and supply something like
				// `--name ""` which will pass the cli lib's requirement
				// check but is still crap
				if *nameOpt == "" {
					panic(errors.New("--name is required"))
				}

				if *roomAliasOpt == "" {
					panic(errors.New("--room is required"))
				} else {
					if ok, roomID = API.Rooms().FindID(*roomAliasOpt); !ok {
						panic(errors.New("could not find room"))
					}
				}

				if *roleNameOpt == "" {
					panic(errors.New("--role is required"))
				} else {
					if ok, roleID = API.RackRoles().FindID(*roleNameOpt); !ok {
						panic(errors.New("could not find rack role"))
					}
				}

				fmt.Println(API.Racks().Create(
					*nameOpt,
					roomID,
					roleID,
					*phaseOpt,
				))
			}
		})
	})

	App.Command("rack", "Work with a single rack", func(cmd *cli.Cmd) {
		var rackID uuid.UUID

		idArg := cmd.StringArg(
			"UUID",
			"",
			"The UUID of the rack. Short UUIDs are *not* accepted, unless you are a Conch sysadmin",
		)

		cmd.Spec = "UUID"

		cmd.Before = func() {

			// BUG(sungo) GetAll() is locked to sysadmin permissions currently.
			// That prevents us from being able to get a full rack list for
			// normal users.
			if IsSysAdmin() {
				var ok bool

				if ok, rackID = API.Racks().FindID(*idArg); !ok {
					panic(errors.New("could not find the rack"))
				}
			} else {
				var err error
				rackID, err = uuid.FromString(*idArg)
				if err != nil {
					panic(err)
				}
			}
		}

		cmd.Command("get", "Get a single rack", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Racks().Get(rackID))
			}
		})

		cmd.Command("update", "Update information about a single rack", func(cmd *cli.Cmd) {
			var (
				nameOpt      = cmd.StringOpt("name", "", "Name of the rack")
				roomAliasOpt = cmd.StringOpt("room", "", "Alias of the datacenter room")
				roleNameOpt  = cmd.StringOpt("role", "", "Name of the role")
				phaseOpt     = cmd.StringOpt("phase", "", "Phase for the rack")

				serialNumberOpt = cmd.StringOpt("serial-number", "", "Serial number of the rack")
				clearSerialOpt  = cmd.BoolOpt("clear-serial-number", false, "Delete the serial number. Overrides --serial-number")

				assetTagOpt      = cmd.StringOpt("asset-tag", "", "Asset Tag of the rack")
				clearAssetTagOpt = cmd.BoolOpt("clear-asset-tag", false, "Delete the asset tag. Overrides --asset-tag")
			)

			cmd.Action = func() {
				var (
					roomID   uuid.UUID
					roleID   uuid.UUID
					ok       bool
					serial   *string
					assetTag *string
				)

				if *roomAliasOpt != "" {
					if ok, roomID = API.Rooms().FindID(*roomAliasOpt); !ok {
						panic(errors.New("could not find room"))
					}
				}
				if *roleNameOpt != "" {
					if ok, roleID = API.RackRoles().FindID(*roleNameOpt); !ok {
						panic(errors.New("could not find rack role"))
					}
				}

				var empty = ""

				if *clearSerialOpt {
					serial = nil
				} else if *serialNumberOpt != "" {
					serial = serialNumberOpt
				} else {
					serial = &empty
				}

				if *clearAssetTagOpt {
					assetTag = nil
				} else if *assetTagOpt != "" {
					assetTag = assetTagOpt
				} else {
					assetTag = &empty
				}

				fmt.Println(API.Racks().Update(
					rackID,
					*nameOpt,
					roomID,
					roleID,
					*phaseOpt,
					serial,
					assetTag,
				))

			}
		})

		cmd.Command("delete rm", "Delete a rack", func(cmd *cli.Cmd) {
			cmd.Before = RequireSysAdmin
			cmd.Action = func() {
				API.Racks().Delete(rackID)
				// BUG(sungo): for sysadmins or GLOBAL admins, this is a lot of
				// data. Maybe should find a better return value here that is
				// still somehow informative
				fmt.Println(API.Racks().GetAll())
			}
		})

		cmd.Command("layout", "The layout of the rack", func(cmd *cli.Cmd) {
			cmd.Command("get", "Get the layout of a rack", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					fmt.Println(API.Racks().Layouts(rackID))
				}
			})

			cmd.Command("export", "Export the layout of the rack as JSON", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					fmt.Println(API.Racks().Layouts(rackID).Export())
				}
			})

			cmd.Command("import", "Import the layout of this rack (using the same format as 'export')", func(cmd *cli.Cmd) {
				var (
					filePathArg  = cmd.StringArg("FILE", "-", "Path to a JSON file that defines the layout. '-' indicates STDIN")
					overwriteOpt = cmd.BoolOpt("overwrite", false, "If the rack has an existing layout, *overwrite* it. This is a destructive action")
				)
				cmd.Action = func() {
					layout := API.Racks().Layouts(rackID)
					if len(layout) > 0 {
						if !*overwriteOpt {
							panic("rack already has a layout. use --overwrite to force")
						}
					}

					var b []byte
					var err error
					if *filePathArg == "-" {
						b, err = ioutil.ReadAll(os.Stdin)
					} else {
						b, err = ioutil.ReadFile(*filePathArg)
					}
					if err != nil {
						panic(err)
					}

					fmt.Println(API.Racks().ImportLayout(rackID, b))
				}
			})
		})

		cmd.Command("assign", "Assign devices to rack slots, using the `--json` output from 'assignments'", func(cmd *cli.Cmd) {
			filePathArg := cmd.StringArg("FILE", "-", "Path to a JSON file to use as the data source. '-' indicates STDIN")
			cmd.Action = func() {

				var b []byte
				var err error
				if *filePathArg == "-" {
					b, err = ioutil.ReadAll(os.Stdin)
				} else {
					b, err = ioutil.ReadFile(*filePathArg)
				}
				if err != nil {
					panic(err)
				}

				fmt.Println(API.Racks().ImportAssignments(rackID, b))

			}
		})

		cmd.Command("assignments", "The devices assigned to the rack", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println(API.Racks().Assignments(rackID))
			}
		})
	})

}
