package main

import (
	"flag"
	"io"
	"log/slog"
	"os"
	"ozon1/internal/commands"
	"ozon1/internal/repository"
	"ozon1/internal/usecase"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)


type Config struct {
	PathToOrders 			string `yaml:"ordersPath"`
	PathToReturnedOrders	string `yaml:"returnPath"`
}

func LoadConfig() *Config {

	path := flag.String("path", "/home/anton/ozon1/config/config.yaml", "path to config")

	flag.Parse()

	file, err := os.Open(*path)

	if err != nil {
		slog.Error("failed to open the file", slog.String("err", err.Error()), slog.String("file", *path))
		return nil
	}

	bytes, err := io.ReadAll(file)

	if err != nil {
		slog.Error("failed to read the file", slog.String("err", err.Error()))
		return nil
	}

	var cfg Config

	err = yaml.Unmarshal(bytes, &cfg)

	if err != nil {
		slog.Error("failed to decode the file", slog.String("err", err.Error()))
		return nil
	}

	return &cfg
}


var rootCmd = &cobra.Command{
	Use:   "pvz",
	Short: "PVZ Order Management",
}

func main() {

	cfg := LoadConfig()

	repo := repository.NewRepository(slog.Default(), cfg.PathToOrders, cfg.PathToReturnedOrders)
	userOperator := usecase.NewUserOperator(repo)
	rootCmd.AddCommand(commands.AcceptOrderFromDelivery(userOperator))
	rootCmd.AddCommand(commands.AcceptReturn(userOperator))
	rootCmd.AddCommand(commands.GetReturnList(userOperator))
	rootCmd.AddCommand(commands.ReturnAllOrders(userOperator))
	rootCmd.AddCommand(commands.ReturnOrderToClient(userOperator))
	rootCmd.AddCommand(commands.ReturnOrderToDelivery(userOperator))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}