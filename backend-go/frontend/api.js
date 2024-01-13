// api.js

// Function to send a POST request
function sendHeartbeat() {
    // Send a POST request to the 'heartbeatz' endpoint
    $.ajax({
        type: 'POST',
        url: '../heartbeatz',  // Adjust the URL accordingly
        data: {
            // You can include any data you want to send with the request
        },
        xhrFields: {
            withCredentials: true
        },
        success: function (response) {
        },
        error: function (error) {
            // Handle the error if needed
            console.error('Error sending heartbeat:', error);
        }
    });
}

// Function to send a GET request to the 'channels' endpoint
function fetchChannels() {
    // Send a GET request to the 'channels' endpoint
    $.ajax({
        type: 'GET',
        url: '../channels',  // Adjust the URL accordingly
        xhrFields: {
            withCredentials: true
        },
        success: function (response) {
            console.log(response);
            // Update the channel list in the HTML based on the response
            updateChannelList(response);
        },
        error: function (error) {
            // Handle the error if needed
            console.error('Error fetching channels:', error);
        }
    });
}
