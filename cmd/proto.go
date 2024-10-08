package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"github.com/yoogoc/kratos-scaffold/generator"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

func newProtoCmd(args []string) *cobra.Command {
	proto := generator.NewProto(settings)
	var protoCmd = &cobra.Command{
		Use:   "proto [NAME]",
		Short: "proto generate req and res model, crud service and xx.pb.go,xx_grpc.pb.go,[xx_http.pb.go].",
		Long: `proto generate req and res model, crud service and xx.pb.go,xx_grpc.pb.go,[xx_http.pb.go].
kratos-scaffold user -n user id:int64:eq,in name:string:contains age:int32:gte,lte --grpc --http`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProto(proto, args)
		},
	}

	addProtoFlags(protoCmd, proto, args)

	return protoCmd
}

func addProtoFlags(protoCmd *cobra.Command, proto *generator.Proto, args []string) {
	flags := protoCmd.PersistentFlags()
	flags.ParseErrorsWhitelist.UnknownFlags = true
	flags.BoolVarP(&proto.GenHttp, "http", "", false, "generate xx.http.pb.go")
	err := flags.Parse(args)
	if err != nil {
		panic(err)
	}
}

func runProto(proto *generator.Proto, args []string) error {
	modelName := args[0]

	proto.Name = util.Singular(strcase.ToCamel(modelName))
	if fs, err := field.ParseFields(args[1:]); err != nil {
		return err
	} else {
		proto.Fields = fs
	}

	err := proto.Generate()
	return err
}
