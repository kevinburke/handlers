http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.7.0/rules_go-0.7.0.tar.gz",
    sha256 = "91fca9cf860a1476abdc185a5f675b641b60d3acf0596679a27b580af60bf19c",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")
go_rules_dependencies()
go_register_toolchains()

go_repository(
    name = "com_github_satori_go_uuid",
    importpath = "github.com/satori/go.uuid",
    commit = "5bf94b69c6b68ee1b541973bb8e1144db23a194b",
)

go_repository(
    name = "com_github_kevinburke_rest",
    importpath = "github.com/kevinburke/rest",
    commit = "5a70172425b1704eedc967b5abf0610866bd26a1",
)

go_repository(
    name = "com_github_mattn_go_colorable",
    importpath = "github.com/mattn/go-colorable",
    commit = "3fa8c76f9daed4067e4a806fb7e4dc86455c6d6a",
)

go_repository(
    name = "com_github_mattn_go_isatty",
    importpath = "github.com/mattn/go-isatty",
    commit = "fc9e8d8ef48496124e79ae0df75490096eccf6fe",
)

go_repository(
    name = "com_github_go_stack_stack",
    importpath = "github.com/go-stack/stack",
    commit = "54be5f394ed2c3e19dac9134a40a95ba5a017f7b",
)

go_repository(
    name = "com_github_inconshreveable_log15",
    importpath = "github.com/inconshreveable/log15",
    commit = "74a0988b5f804e8ce9ff74fca4f16980776dff29",
)
