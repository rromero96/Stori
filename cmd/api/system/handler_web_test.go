package system_test

/*  func TestHTTPHandler_GetHTMLInfoV1_success(t *testing.T) {
	processTransaction := system.MockHTMLProcessTransactions([]byte{}, nil)
	getHTMLInfoV1 := system.GetHTMLInfoV1(processTransaction)

	ctx, w := context.Background(), httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/test", strings.NewReader(""))

	got := getHTMLInfoV1(w, r)

	assert.Nil(t, got)
}

func TestHTTPHandler_GetHTMLInfoV1_fails(t *testing.T) {
	processTransaction := system.MockHTMLProcessTransactions([]byte{}, system.ErrCantGetTransactionInfo)
	getHTMLInfoV1 := system.GetHTMLInfoV1(processTransaction)

	ctx, w := context.Background(), httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/test", strings.NewReader(""))

	want := web.NewError(http.StatusInternalServerError, system.CantGetInfo)
	got := getHTMLInfoV1(w, r)

	assert.Equal(t, got, want)
}

*/
