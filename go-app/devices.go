package main

// Device represents hardware device
type Device struct {
	// Universally unique identifier
	UUID string `json:"UUID"`

	// Mac address
	Mac string `json:"mac"`

	// Firmware version
	Firmware string `json:"firmware"`
}

// devices returns pseudo connected devices.
func devices() []Device {
	return []Device{
		{UUID: "b0e42fe7-31a5-4894-a441-007e5256afea", Mac: "5F-33-CC-1F-43-82", Firmware: "2.1.6"},
		{UUID: "0c3242f5-ae1f-4e0c-a31b-5ec93825b3e7", Mac: "EF-2B-C4-F5-D6-34", Firmware: "2.1.5"},
		{UUID: "b16d0b53-14f1-4c11-8e29-b9fcef167c26", Mac: "62-46-13-B7-B3-A1", Firmware: "3.0.0"},
		{UUID: "51bb1937-e005-4327-a3bd-9f32dcf00db8", Mac: "96-A8-DE-5B-77-14", Firmware: "1.0.1"},
		{UUID: "e0a1d085-dce5-48db-a794-35640113fa67", Mac: "7E-3B-62-A6-09-12", Firmware: "3.5.6"},
		{UUID: "f47ac10b-58cc-4372-a567-0e02b2c3d479", Mac: "A2-8E-F1-4D-9C-7B", Firmware: "4.2.1"},
		{UUID: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", Mac: "D8-BB-2C-E6-A1-F3", Firmware: "1.8.9"},
		{UUID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Mac: "B4-E9-B0-F2-8A-5C", Firmware: "5.1.2"},
		{UUID: "7c9e6679-7425-40de-944b-e07fc1f90ae7", Mac: "C6-4A-3D-9E-7F-B2", Firmware: "2.9.4"},
		{UUID: "a1b2c3d4-e5f6-7890-1234-567890abcdef", Mac: "F0-9F-C2-3A-8B-4E", Firmware: "6.0.3"},
		{UUID: "12345678-1234-5678-1234-123456789abc", Mac: "8C-16-45-AC-D2-B7", Firmware: "3.7.8"},
		{UUID: "98765432-4321-8765-4321-cba987654321", Mac: "E4-A4-71-CC-8F-D9", Firmware: "4.5.1"},
		{UUID: "11111111-2222-3333-4444-555555555555", Mac: "9A-2F-B8-6C-E3-A5", Firmware: "1.4.7"},
		{UUID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", Mac: "7D-C8-91-A2-F4-B6", Firmware: "7.2.0"},
		{UUID: "fedcba98-7654-3210-fedc-ba9876543210", Mac: "B8-27-EB-D4-A1-9C", Firmware: "2.3.5"},
	}
}
