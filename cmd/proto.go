package cmd

import (
	"github.com/YoogoC/kratos-scaffold/generator"
	"github.com/YoogoC/kratos-scaffold/pkg/field"
	"github.com/YoogoC/kratos-scaffold/pkg/util"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

func newProtoCmd() *cobra.Command {
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

	addProtoFlags(protoCmd, proto)
	return protoCmd
}

func addProtoFlags(protoCmd *cobra.Command, proto *generator.Proto) {
	protoCmd.PersistentFlags().BoolVarP(&proto.GenHttp, "http", "", true, "generate xx.http.pb.go")
	protoCmd.PersistentFlags().BoolVarP(&proto.GenHttp, "grpc", "", true, "generate xx.grpc.pb.go")
}

func runProto(proto *generator.Proto, args []string) error {
	modelName := args[0]

	proto.Name = util.Singular(strcase.ToCamel(modelName))
	proto.Fields = field.ParseFields(args[1:])

	err := proto.Generate()
	return err
}
