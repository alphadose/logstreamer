package main

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "logstreamer",
	Short: "LogStreamer CLI",
	Long:  `LogStreamer CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		// if len(args) == 0 {
		// 	cmd.Usage()
		// 	utils.GracefulExit()
		// }
		// filepath := args[0]
		// var contents, err = ioutil.ReadFile(filepath)
		// if err != nil {
		// 	fmt.Println(err)
		// 	utils.GracefulExit()
		// }

		// aPIDefinition := openapiclient.NewAPIDefinition()
		// if err := json.Unmarshal(contents, aPIDefinition); err != nil {
		// 	fmt.Println(err)
		// 	utils.GracefulExit()
		// }

		// aPIDefinition.SetName("TestAPI")
		// aPIDefinition.SetUseKeyless(true)
		// aPIDefinition.SetActive(true)

		// proxy := *openapiclient.NewAPIDefinitionProxy()
		// proxy.SetListenPath("/test")

		// aPIDefinition.SetProxy(proxy)

		// utils.Out(utils.GetDefaultClient().APIsApi.CreateApi(utils.GetAuthCtx()).APIDefinition(*aPIDefinition).Execute())
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.logstreamer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
