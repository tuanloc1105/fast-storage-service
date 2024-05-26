package implement

import (
	"context"
	"fast-storage-go-service/utils"
	"fmt"
)

func generateSecret(ctx context.Context) (string, error) {
	if generateSecretShellStd, _, generateSecretError := utils.Shellout(ctx, "java -jar additional_source_code/two-factor-auth.jar \"GENERATE_BASE32_SECRET\""); generateSecretError != nil {
		return "", generateSecretError
	} else {
		return generateSecretShellStd, nil
	}
}

func GenerateTotp(ctx context.Context, otpSecretKey string) (string, error) {
	otpStdOut, _, otpError := utils.Shellout(
		ctx,
		fmt.Sprintf(
			"java -jar additional_source_code/two-factor-auth.jar \"GENERATE_CURRENT_OTP\" \"%s\"",
			otpSecretKey,
		),
	)
	return otpStdOut, otpError
}
