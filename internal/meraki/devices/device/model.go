package device

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*

"/devices/{serial}": {
      "get": {
        "description": "Return a single device",
        "operationId": "getDevice",
        "parameters": [
          {
            "name": "serial",
            "in": "path",
            "description": "Serial",
            "schema": {
              "type": "string"
            },
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Successful operation",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "name": {
                      "type": "string",
                      "description": "Name of the device"
                    },
                    "lat": {
                      "type": "number",
                      "format": "float",
                      "description": "Latitude of the device"
                    },
                    "lng": {
                      "type": "number",
                      "format": "float",
                      "description": "Longitude of the device"
                    },
                    "address": {
                      "type": "string",
                      "description": "Physical address of the device"
                    },
                    "notes": {
                      "type": "string",
                      "description": "Notes for the device, limited to 255 characters"
                    },
                    "tags": {
                      "type": "array",
                      "items": {
                        "type": "string"
                      },
                      "description": "List of tags assigned to the device"
                    },
                    "networkId": {
                      "type": "string",
                      "description": "ID of the network the device belongs to"
                    },
                    "serial": {
                      "type": "string",
                      "description": "Serial number of the device"
                    },
                    "model": {
                      "type": "string",
                      "description": "Model of the device"
                    },
                    "mac": {
                      "type": "string",
                      "description": "MAC address of the device"
                    },
                    "lanIp": {
                      "type": "string",
                      "description": "LAN IP address of the device"
                    },
                    "firmware": {
                      "type": "string",
                      "description": "Firmware version of the device"
                    },
                    "floorPlanId": {
                      "type": "string",
                      "description": "The floor plan to associate to this device. null disassociates the device from the floorplan."
                    },
                    "details": {
                      "type": "array",
                      "items": {
                        "type": "object",
                        "properties": {
                          "name": {
                            "type": "string",
                            "description": "Additional property name"
                          },
                          "value": {
                            "type": "string",
                            "description": "Additional property value"
                          }
                        }
                      },
                      "description": "Additional device information"
                    },
                    "beaconIdParams": {
                      "type": "object",
                      "properties": {
                        "uuid": {
                          "type": "string",
                          "description": "The UUID to be used in the beacon identifier"
                        },
                        "major": {
                          "type": "integer",
                          "description": "The major number to be used in the beacon identifier"
                        },
                        "minor": {
                          "type": "integer",
                          "description": "The minor number to be used in the beacon identifier"
                        }
                      },
                      "description": "Beacon Id parameters with an identifier and major and minor versions"
                    }
                  }
                },
                "example": {
                  "name": "My AP",
                  "lat": 37.4180951010362,
                  "lng": -122.098531723022,
                  "address": "1600 Pennsylvania Ave",
                  "notes": "My AP's note",
                  "tags": [
                    " recently-added "
                  ],
                  "networkId": "N_24329156",
                  "serial": "Q234-ABCD-5678",
                  "model": "MR34",
                  "mac": "00:11:22:33:44:55",
                  "lanIp": "1.2.3.4",
                  "firmware": "wireless-25-14",
                  "floorPlanId": "g_2176982374",
                  "details": [
                    {
                      "name": "Catalyst serial",
                      "value": "123ABC"
                    }
                  ],
                  "beaconIdParams": {
                    "uuid": "00000000-0000-0000-0000-000000000000",
                    "major": 5,
                    "minor": 3
                  }
                }
              }
            }
          }
        },
        "summary": "Return a single device",
        "tags": [
          "devices",
          "configure"
        ]
      },
      "put": {
        "description": "Update the attributes of a device",
        "operationId": "updateDevice",
        "parameters": [
          {
            "name": "serial",
            "in": "path",
            "description": "Serial",
            "schema": {
              "type": "string"
            },
            "required": true
          }
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "name": {
                    "type": "string",
                    "description": "The name of a device"
                  },
                  "tags": {
                    "type": "array",
                    "items": {
                      "type": "string"
                    },
                    "description": "The list of tags of a device"
                  },
                  "lat": {
                    "type": "number",
                    "format": "float",
                    "description": "The latitude of a device"
                  },
                  "lng": {
                    "type": "number",
                    "format": "float",
                    "description": "The longitude of a device"
                  },
                  "address": {
                    "type": "string",
                    "description": "The address of a device"
                  },
                  "notes": {
                    "type": "string",
                    "description": "The notes for the device. String. Limited to 255 characters."
                  },
                  "moveMapMarker": {
                    "type": "boolean",
                    "description": "Whether or not to set the latitude and longitude of a device based on the new address. Only applies when lat and lng are not specified."
                  },
                  "switchProfileId": {
                    "type": "string",
                    "description": "The ID of a switch template to bind to the device (for available switch templates, see the 'Switch Templates' endpoint). Use null to unbind the switch device from the current profile. For a device to be bindable to a switch template, it must (1) be a switch, and (2) belong to a network that is bound to a configuration template."
                  },
                  "floorPlanId": {
                    "type": "string",
                    "description": "The floor plan to associate to this device. null disassociates the device from the floorplan."
                  }
                },
                "example": {
                  "name": "My AP",
                  "tags": [
                    " recently-added "
                  ],
                  "lat": 37.4180951010362,
                  "lng": -122.098531723022,
                  "address": "1600 Pennsylvania Ave",
                  "notes": "My AP's note",
                  "moveMapMarker": true,
                  "switchProfileId": "1234",
                  "floorPlanId": "g_2176982374"
                }
              }
            }
          },
          "required": false
        },
        "responses": {
          "200": {
            "description": "Successful operation",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "name": {
                      "type": "string",
                      "description": "Name of the device"
                    },
                    "lat": {
                      "type": "number",
                      "format": "float",
                      "description": "Latitude of the device"
                    },
                    "lng": {
                      "type": "number",
                      "format": "float",
                      "description": "Longitude of the device"
                    },
                    "address": {
                      "type": "string",
                      "description": "Physical address of the device"
                    },
                    "notes": {
                      "type": "string",
                      "description": "Notes for the device, limited to 255 characters"
                    },
                    "tags": {
                      "type": "array",
                      "items": {
                        "type": "string"
                      },
                      "description": "List of tags assigned to the device"
                    },
                    "networkId": {
                      "type": "string",
                      "description": "ID of the network the device belongs to"
                    },
                    "serial": {
                      "type": "string",
                      "description": "Serial number of the device"
                    },
                    "model": {
                      "type": "string",
                      "description": "Model of the device"
                    },
                    "mac": {
                      "type": "string",
                      "description": "MAC address of the device"
                    },
                    "lanIp": {
                      "type": "string",
                      "description": "LAN IP address of the device"
                    },
                    "firmware": {
                      "type": "string",
                      "description": "Firmware version of the device"
                    },
                    "floorPlanId": {
                      "type": "string",
                      "description": "The floor plan to associate to this device. null disassociates the device from the floorplan."
                    },
                    "details": {
                      "type": "array",
                      "items": {
                        "type": "object",
                        "properties": {
                          "name": {
                            "type": "string",
                            "description": "Additional property name"
                          },
                          "value": {
                            "type": "string",
                            "description": "Additional property value"
                          }
                        }
                      },
                      "description": "Additional device information"
                    },
                    "beaconIdParams": {
                      "type": "object",
                      "properties": {
                        "uuid": {
                          "type": "string",
                          "description": "The UUID to be used in the beacon identifier"
                        },
                        "major": {
                          "type": "integer",
                          "description": "The major number to be used in the beacon identifier"
                        },
                        "minor": {
                          "type": "integer",
                          "description": "The minor number to be used in the beacon identifier"
                        }
                      },
                      "description": "Beacon Id parameters with an identifier and major and minor versions"
                    }
                  }
                },
                "example": {
                  "name": "My AP",
                  "lat": 37.4180951010362,
                  "lng": -122.098531723022,
                  "address": "1600 Pennsylvania Ave",
                  "notes": "My AP's note",
                  "tags": [
                    " recently-added "
                  ],
                  "networkId": "N_24329156",
                  "serial": "Q234-ABCD-5678",
                  "model": "MR34",
                  "mac": "00:11:22:33:44:55",
                  "lanIp": "1.2.3.4",
                  "firmware": "wireless-25-14",
                  "floorPlanId": "g_2176982374",
                  "details": [
                    {
                      "name": "Catalyst serial",
                      "value": "123ABC"
                    }
                  ],
                  "beaconIdParams": {
                    "uuid": "00000000-0000-0000-0000-000000000000",
                    "major": 5,
                    "minor": 3
                  }
                }
              }
            }
          }
        },
        "summary": "Update the attributes of a device",
        "tags": [
          "devices",
          "configure"
        ]
      }
    },
*/

// ResourceModel represents the device resource in Terraform.
type ResourceModel struct {
	Id              types.String  `tfsdk:"id" json:"id"`
	Serial          types.String  `tfsdk:"serial" json:"serial"`
	Name            types.String  `tfsdk:"name" json:"name"`
	Mac             types.String  `tfsdk:"mac" json:"mac"`
	Model           types.String  `tfsdk:"model" json:"model"`
	Tags            types.List    `tfsdk:"tags" json:"tags"`
	Details         types.List    `tfsdk:"details" json:"details"`
	LanIp           types.String  `tfsdk:"lan_ip" json:"lanIp"`
	Firmware        types.String  `tfsdk:"firmware" json:"firmware"`
	Lat             types.Float64 `tfsdk:"lat" json:"lat"`
	Lng             types.Float64 `tfsdk:"lng" json:"lng"`
	Address         types.String  `tfsdk:"address" json:"address"`
	Notes           types.String  `tfsdk:"notes" json:"notes"`
	Url             types.String  `tfsdk:"url" json:"url"`
	FloorPlanId     types.String  `tfsdk:"floor_plan_id" json:"floorPlanId"`
	NetworkId       types.String  `tfsdk:"network_id" json:"networkId"`
	BeaconIdParams  types.Object  `tfsdk:"beacon_id_params" json:"beaconIdParams"`
	SwitchProfileId types.String  `tfsdk:"switch_profile_id" json:"switchProfileId"`
	MoveMapMarker   types.Bool    `tfsdk:"move_map_marker" json:"moveMapMarker"`
}

// BeaconIdParamsModel represents the beacon ID parameters of a device.
type BeaconIdParamsModel struct {
	Uuid  types.String `tfsdk:"uuid" json:"uuid"`
	Major types.Int64  `tfsdk:"major" json:"major"`
	Minor types.Int64  `tfsdk:"minor" json:"minor"`
}

// BeaconIdParamsType represents the beacon ID parameters of a device.
var BeaconIdParamsType = map[string]attr.Type{
	"uuid":      types.StringType,
	"major":     types.Int64Type,
	"minor":     types.Int64Type,
	"beacon_id": types.StringType,
	"proximity": types.StringType,
}

// DetailsModel represents additional details of a device.
type DetailsModel struct {
	Name  types.String `tfsdk:"name" json:"name"`
	Value types.String `tfsdk:"value" json:"value"`
}

var DetailsType = map[string]attr.Type{
	"name":  types.StringType,
	"value": types.StringType,
}
