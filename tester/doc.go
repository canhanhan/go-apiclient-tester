/*
Package tester implements a library that makes testing HTTP based API clients easier.

Example test:
	func setup() *tester.Tester {
		return tester.NewTester(map[string]tester.TestCategory{
			"users": tester.TestCategory{
				Scenarios: map[string]tester.TestScenario{
					"create": TestScenario{
						Request: TestRequest{
							Method: "POST",
							Path:   "/user",
							Headers: map[string]string{"Content-Type": "application/json"},
							Body: strings.NewReader("{\"username\":\"user1\"}"),
						},
						Response: TestResponse{
							Code: 200,
							Headers: map[string]string{"Content-Type": "application/json"},
							Body: strings.NewReader("{\"id\":1,\"username\":\"user1\"}"),
						},
					},
				},
			}
		})
	}

	func TestUserCreate(t *testing.T) {
		api := setup()
		defer api.Close()
		api.Setup(t, "user", "create")
		c := NewClient(api.URL)

		user, err := c.CreateUser("test")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "user1", user.Name)
	}
*/
package tester
