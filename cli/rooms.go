package cli

import (
	"errors"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/kosh/v3/conch"
	"github.com/joyent/kosh/v3/conch/types"
)

func roomsCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display func(interface{})

	cmd.Before = func() {
		requireSysAdmin(config)()
		conch = config.ConchClient()
		display = config.Renderer()
	}

	cmd.Command("get", "Get a list of all rooms", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			display(conch.GetAllRooms())
		}
	})

	cmd.Command("create", "Create a single room", func(cmd *cli.Cmd) {
		var (
			aliasOpt        = cmd.StringOpt("alias", "", "Alias")
			azOpt           = cmd.StringOpt("az", "", "AZ")
			datacenterIDOpt = cmd.StringOpt("datacenter-id", "", "Datacenter UUID (first segment of UUID accepted)")
			vendorNameOpt   = cmd.StringOpt("vendor-name", "", "Vendor Name")
		)

		cmd.Spec = "--datacenter-id --alias --az [OPTIONS]"
		cmd.Action = func() {
			// The user can be very silly and supply something like
			// '--alias ""' which will pass the cli lib's requirement
			// check but is still crap
			if *aliasOpt == "" {
				fatal(errors.New("--alias is required"))
			}
			if *azOpt == "" {
				fatal(errors.New("--az is required"))
			}
			if *datacenterIDOpt == "" {
				fatal(errors.New("--datacenter-id is required"))
			}

			datacenter := conch.GetDatacenterByName(*datacenterIDOpt)
			if (datacenter == types.Datacenter{}) {
				fatal(errors.New("could not find the datacenter"))
			}

			conch.CreateRoom(types.DatacenterRoomCreate{
				DatacenterID: datacenter.ID,
				Az:           types.NonEmptyString(*azOpt),
				Alias:        types.MojoStandardPlaceholder(*aliasOpt),
				VendorName:   types.MojoRelaxedPlaceholder(*vendorNameOpt),
			})
		}
	})
}

func roomCmd(cmd *cli.Cmd) {
	var conch *conch.Client
	var display func(interface{})
	var room types.DatacenterRoomDetailed

	aliasArg := cmd.StringArg(
		"ALIAS",
		"",
		"The unique alias of the datacenter room",
	)

	cmd.Spec = "ALIAS"

	cmd.Before = func() {
		requireSysAdmin(config)()

		conch = config.ConchClient()
		display = config.Renderer()

		room = conch.GetRoomByAlias(*aliasArg)
		if (room == types.DatacenterRoomDetailed{}) {
			fatal(errors.New("could not find the room"))
		}
	}

	cmd.Command("get", "Information about a single room", func(cmd *cli.Cmd) {
		cmd.Action = func() { display(room) }
	})

	cmd.Command("update", "Update information about a single room", func(cmd *cli.Cmd) {
		var (
			aliasOpt        = cmd.StringOpt("alias", "", "Alias")
			azOpt           = cmd.StringOpt("az", "", "AZ")
			datacenterIDOpt = cmd.StringOpt("datacenter-id", "", "Datacenter UUID (first segment of UUID accepted)")
			vendorNameOpt   = cmd.StringOpt("vendor-name", "", "Vendor Name")
		)

		cmd.Action = func() {
			dc := conch.GetDatacenterByName(*datacenterIDOpt)
			if (dc == types.Datacenter{}) {
				fatal(errors.New("could not find the datacenter"))
			}

			conch.UpdateRoom(room.ID, types.DatacenterRoomUpdate{
				DatacenterID: dc.ID,
				Az:           types.NonEmptyString(*azOpt),
				Alias:        types.MojoStandardPlaceholder(*aliasOpt),
				VendorName:   types.MojoRelaxedPlaceholder(*vendorNameOpt),
			})
			display(conch.GetRoomByID(room.ID))
		}
	})

	cmd.Command("delete", "Delete a single room", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			conch.DeleteRoom(room.ID)
			display(conch.GetAllRooms())
		}
	})

	cmd.Command("racks", "View the racks assigned to a single room", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			display(conch.GetAllRoomRacks(room.ID))
		}
	})
}
