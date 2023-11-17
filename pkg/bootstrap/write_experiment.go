package bootstrap

import (
	"embed"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"os"
)

//go:embed example.tf
var terraformFile embed.FS

func ReadHCLFile() {
	// Read the embedded Terraform file.
	content, err := terraformFile.ReadFile("example.tf")
	if err != nil {
		fmt.Println("Error reading the embedded file:", err)
		return
	}

	// Parse the HCL.
	file, diag := hclwrite.ParseConfig(content, "example.tf", hcl.InitialPos)
	if diag.HasErrors() {
		fmt.Println("Error parsing HCL:", diag)
		return
	}

	// Modify the HCL data structure as needed.
	// This will depend on your specific requirements.
	rootBody := file.Body()
	//rootBody.GetAttribute("locals").SetAttributeValue("foo", cty.StringVal("bar")
	//rootBody.GetAttribute()
	localsBlock := rootBody.FirstMatchingBlock("locals", nil)
	if localsBlock != nil {
		// Assuming you want to add or update a local variable named "example"
		localsBlock.Body().SetAttributeValue("hei", cty.StringVal("yoooooooooooooooooooo"))
	}

	rootBody.SetAttributeValue("string", cty.StringVal("bar")) // this is overwritten later
	rootBody.AppendNewline()
	rootBody.SetAttributeValue("object", cty.ObjectVal(map[string]cty.Value{
		"foo": cty.StringVal("foo"),
		"bar": cty.NumberIntVal(5),
		"baz": cty.True,
	}))

	// Serialize and write the modified HCL back to a file.
	modifiedContent := file.Bytes()
	err = os.WriteFile("modified_example.tf", modifiedContent, 0644)
	if err != nil {
		fmt.Println("Error writing modified file:", err)
		return
	}

	fmt.Println("The Terraform file has been modified successfully.")
}
