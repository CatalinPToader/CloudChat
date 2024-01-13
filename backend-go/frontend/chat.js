// chat.js

// Function to send a chat message
function sendChatMessage(channelName) {
    // Get the message from the message-input textarea
    const message = $('#message-input').val().trim();

    // If the message is not empty, send a POST request
    if (message !== '') {
        // Send a POST request to the 'channel/$channelName' endpoint
        $.ajax({
            type: 'POST',
            url: `../channel/${channelName}`,  // Use the current channel
            xhrFields: {
                withCredentials: true
            },
            contentType: 'application/json',
            data: JSON.stringify({
                "message": message
            }),
            success: function (response) {
            },
            error: function (error) {
                // Handle the error if needed
                console.error(`Error sending message for channel ${channelName}:`, error);
            }
        });

        // Clear the message-input textarea after sending the message
        $('#message-input').val('');
    }
}

// Event listener for the send button click
$('#send-button').click(function () {
    sendChatMessage(currentChannel);
});

// Event listener for the message-input textarea keypress (Enter key)
$('#message-input').keypress(function (event) {
    if (event.which === 13) { // Enter key
        sendChatMessage(currentChannel);
        event.preventDefault(); // Prevent the default behavior of the Enter key (new line in textarea)
    }
});


// Variable to track the current channel
let currentChannel = '';

// Interval for continuous message fetching
let fetchInterval;

// Function to start continuous message fetching
function startContinuousFetch(channelName) {
    // Clear any existing interval
    clearInterval(fetchInterval);

    // Update the current channel
    currentChannel = channelName;

    // Fetch messages immediately
    fetchMessages();

    // Set up an interval for continuous fetching every second
    fetchInterval = setInterval(fetchMessages, 1000);
}

// Function to stop continuous message fetching
function stopContinuousFetch() {
    // Clear the interval
    clearInterval(fetchInterval);
}

// Function to send a fetch request for messages
function fetchMessages() {
    // Send a GET request to the 'channel/$channelName' endpoint
    $.ajax({
        type: 'GET',
        url: `../channel/${currentChannel}`,  // Use the current channel
        xhrFields: {
            withCredentials: true
        },
        success: function (response) {
            // Update the chat messages div based on the response
            updateChatMessages(response);
        },
        error: function (error) {
            // Handle the error if needed
            console.error(`Error fetching messages for channel ${currentChannel}:`, error);
        }
    });
}

// Function to update the chat messages div with fetched messages
function updateChatMessages(response) {
    // Assuming the response is in JSON format
    const messageList = response && response.message_list ? response.message_list : [];

    // Clear the existing chat messages div only if there are new messages
    if (messageList.length > 0) {
        $('#chat-messages').empty();
    }

    // Append the messages to the chat messages div, avoiding duplicates
    messageList.forEach(message => {
        const timestamp = new Date(message.timestamp * 1000);
        const formattedTimestamp = `${timestamp.getFullYear()}/${(timestamp.getMonth() + 1).toString().padStart(2, '0')}/${timestamp.getDate().toString().padStart(2, '0')} ${timestamp.getHours().toString().padStart(2, '0')}:${timestamp.getMinutes().toString().padStart(2, '0')}`;
        const formattedMessage = `${formattedTimestamp} ${message.username}: ${message.message}`;

        // Check if the message is not a duplicate before appending
        if (!$(`#chat-messages:contains('${formattedMessage}')`).length) {
            $('#chat-messages').append(`<div>${formattedMessage}</div>`);
        }
    });
}   

// Function to change the "Chat Room" text
function changeChatRoomText(channelName) {
    $('#chat-title').text(`Chat Room - ${channelName}`);
}

// Function to clear the chat messages div
function clearChatMessages() {
    $('#chat-messages').empty();
}
