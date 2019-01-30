package main

// Test getSubmissions:
// - Auth Token is correctly set
// - Correct target is the URL
// - Correctly unmarshalled response is returned
// - Error is returned if the auth fails
// - Error is returned if the url is wrong

// Test checkSubmissions:
// - Check that unknown submissions are detected
// - Check that changed submissions are detected
// - Check that the notification is invoked