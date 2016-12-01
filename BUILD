load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_prefix",
    "go_library",
    "go_binary",
    "go_test",
)

# Define an import prefix
go_prefix("github.com/cloudendpoints/endpoints-tools")

go_library(
    name = "vendor/deploy",
    srcs = [
        "deploy/service.go",
        "deploy/default_service_config.go",
    ],
    go_prefix = "//:go_prefix",
    deps = [
        "@golang_oauth2//:oauth2/google",
        "@golang_oauth2//:oauth2/jwt",
        "@github_com_x_net//:context",
        "@google_api_go_client//:servicemanagement/v1",
        "@github_com_golang_protobuf//:jsonpb",
        "@github_com_golang_protobuf//:proto",
        "@go_yaml//:yaml.v2",
    ],
)

go_library(
    name = "vendor/cli",
    srcs = glob([
        "cli/*.go",
    ]),
    go_prefix = "//:go_prefix",
    deps = [
        ":vendor/deploy",
        "@google_api_go_client//:logging/v2beta1",
        "@golang_oauth2//:oauth2/google",
        "@github_com_x_net//:context",
        "@github_com_spf13_cobra//:cobra",
        "@github_com_kubernetes_client_go//:kubernetes",
        "@github_com_kubernetes_client_go//:pkg/labels",
        "@github_com_kubernetes_client_go//:tools/clientcmd",
        "@github_com_kubernetes_client_go//:pkg/api",
        "@github_com_kubernetes_client_go//:pkg/api/v1",
        "@github_com_kubernetes_client_go//:pkg/util/intstr",
        "@github_com_kubernetes_client_go//:pkg/apis/extensions/v1beta1",
    ],
)

go_binary(
    name = "espcli",
    srcs = [
        "espcli.go",
    ],
    go_prefix = "//:go_prefix",
    deps = [
        ":vendor/cli",
        "@github_com_spf13_cobra//:cobra",
    ],
)

go_binary(
    name = "service_config",
    srcs = [
        "service_config.go",
    ],
    go_prefix = "//:go_prefix",
    deps = [
        ":vendor/deploy",
        "@google_api_go_client//:servicemanagement/v1",
        "@github_com_spf13_cobra//:cobra",
    ],
)
