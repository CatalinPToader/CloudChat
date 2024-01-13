// users.js

// Variable to track the current channel for user fetching
let currentUserChannel = '';

// Interval for continuous user fetching
let userFetchInterval;

// Function to start continuous user fetching
function startContinuousUserFetch(channelName) {
    // Clear any existing interval
    clearInterval(userFetchInterval);

    // Update the current user channel
    currentUserChannel = channelName;

    // Fetch users immediately
    fetchUsers();

    // Set up an interval for continuous user fetching every second
    userFetchInterval = setInterval(fetchUsers, 1000);
}

// Function to stop continuous user fetching
function stopContinuousUserFetch() {
    // Clear the interval
    clearInterval(userFetchInterval);
}

// Function to send a GET request to fetch users for a specific channel
function fetchUsers() {
    // Send a GET request to the 'users/$channelName' endpoint
    $.ajax({
        type: 'GET',
        url: `../users/${currentUserChannel}`,  // Use the current user channel
        xhrFields: {
            withCredentials: true
        },
        success: function (response) {
            // Update the user list in the HTML based on the response
            updateUserList(response);
        },
        error: function (error) {
            // Handle the error if needed
            console.error(`Error fetching users for channel ${currentUserChannel}:`, error);
        }
    });
}

// Function to update the user list in the HTML
function updateUserList(response) {
    // Assuming the response is in JSON format
    const userList = response && response.user_list ? response.user_list : [];

    // Clear the existing user list
    $('#users-list').empty();

    // Append the updated user list to the HTML
    userList.forEach(user => {
        $('#users-list').append(`<li>${user}</li>`);
    });
}

// Function to clear the chat messages div
function clearUserList() {
    $('#users-list').empty();
}