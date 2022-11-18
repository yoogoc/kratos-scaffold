package project_generator

func genBffInternal(name string, appPath string) error {
	// biz,data,service,server,conf
	// orm := "grpc"
	// 1. mkdir biz. gen biz/biz.go
	if err := genBffBiz(appPath); err != nil {
		return err
	}

	// 2. mkdir service. gen service/service.go
	if err := genService(appPath); err != nil {
		return err
	}

	// 3. gen data
	if err := genBffData(name, appPath); err != nil {
		return err
	}

	// 4 mkdir server. gen server
	if err := genServer(name, appPath, true); err != nil {
		return err
	}

	// 5 gen conf
	if err := genConf(name, appPath, confBffProto); err != nil {
		return err
	}
	return nil
}
