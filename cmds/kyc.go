package cmds

type KYCCommand struct {
	CreateKYCService  CreateKYCServiceCommand  `cmd:"" name:"create-kyc-service" help:"create kyc service to contract account"`
	AddControllers    AddControllersCommand    `cmd:"" name:"add-controllers" help:"add controllers to kyc service"`
	RemoveControllers RemoveControllersCommand `cmd:"" name:"remove-controllers" help:"remove controllers from key service"`
	AddCustomers      AddCustomersCommand      `cmd:"" name:"add-customers" help:"add customer status to kyc service"`
	UpdateCustomers   UpdateCustomersCommand   `cmd:"" name:"update-customers" help:"update registered customer status"`
}
