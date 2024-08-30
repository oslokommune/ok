package pkg

/*
var TestCommand = &cobra.Command{
	Use: "test",
	RunE: func(cmd *cobra.Command, args []string) error {
		var v any = make(map[string]any)
		bin, err := os.ReadFile("app-v8.0.2.dependencies.yml")
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(bin, &v); err != nil {
			return err
		}
		fmt.Printf("Type of root: %T\n", v)
		vasmap, ok := v.(map[string]any)
		if !ok {
			fmt.Printf("NOPE\n")
			return fmt.Errorf("not a map")
		}

		for k, vv := range vasmap {
			fmt.Printf("Type of %s: %T\n", k, vv)
			if _, ok := vv.(map[string]any); ok {
				fmt.Printf("YEAH IT MATCHES map any!\n")
			} else {
				fmt.Printf("NOPE\n")
			}
		}

		return nil
	},
}
*/
/*
var SchemaCommand = &cobra.Command{
	Use: "schema dependencies-input schema-output",
	RunE: func(cmd *cobra.Command, args []string) error {

		inputFile, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer inputFile.Close()
		schemaFileName := args[1]
		dec := yaml.NewDecoder(inputFile)
		var dependencies = make(map[string]*DownloadedBoilThingy)
		if err = dec.Decode(&dependencies); err != nil {
			return err
		}

		rootCfg, err := findRootConfig(dependencies)
		if err != nil {
			return err
		}
		allVariables := collectFolderVariables("", rootCfg.Path, rootCfg, dependencies)
		jsonSchema := buildJsonSchemaFromNamespaceVariables(allVariables)
		cmd.Printf("Writing schema file to %s\n", schemaFileName)
		if err := writeJsonSchemaToFile(jsonSchema, schemaFileName); err != nil {
			return err
		}

		return nil
	},
}
*/
func init() {

	/*
		ConfigCommand.AddCommand(SchemaCommand)
		ConfigCommand.AddCommand(TestCommand)
	*/
}
