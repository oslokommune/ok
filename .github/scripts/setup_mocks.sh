function setupMocks() {
    function terraform() {
        printf "mock terraform with args: %s\n" "$@"
    }
    export -f terraform

    function fzf() {
        printf "mock fzf with args: %s\n" "$@"
    }
    export -f fzf

    function gh() {
        printf "mock gh with args: %s\n" "$@"
    }
    export -f gh
}
