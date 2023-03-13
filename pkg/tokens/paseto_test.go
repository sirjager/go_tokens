package tokens

import (
	"testing"
	"time"

	"github.com/sirjager/go_tokens/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoBuilder(t *testing.T) {
	builder, err := NewPasetoBuilder(small_secret_key)
	require.Error(t, err)
	require.Empty(t, builder)

	builder, err = NewPasetoBuilder(valid_secret_key)
	require.NoError(t, err)
	require.NotEmpty(t, builder)

	// Create Token
	data := PayloadData{Data: utils.RandomString(64)}

	hash, payload, err := builder.CreateToken(data, time.Second*10)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, payload.Id)
	require.NotEmpty(t, payload.IssuedAt)
	require.NotEmpty(t, payload.ExpiredAt)
	require.Equal(t, data.Data, payload.Payload.Data)
	// Now verify
	payload, err = builder.VerifyToken(hash)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, payload.Id)
	require.NotEmpty(t, payload.IssuedAt)
	require.NotEmpty(t, payload.ExpiredAt)
	require.NotEmpty(t, payload.Payload.Data)
	require.Equal(t, data.Data, payload.Payload.Data)

	// Verify Token
	// with expired token
	hash, payload, err = builder.CreateToken(data, time.Microsecond)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, payload.Id)
	require.NotEmpty(t, payload.IssuedAt)
	require.NotEmpty(t, payload.ExpiredAt)
	require.Equal(t, data.Data, payload.Payload.Data)
	// Now verify with expired token
	payload, err = builder.VerifyToken(hash)
	require.Error(t, err)
	require.Empty(t, payload)
	require.EqualError(t, ErrExpiredToken, err.Error())
}
