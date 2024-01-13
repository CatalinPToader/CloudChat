// channels.js

// Function to update the channel list in the HTML
function updateChannelList(response) {
    // Assuming the response is in JSON format
    const channelList = response && response.channel_list ? response.channel_list : [];

    // Clear the existing channel list
    $('#channels-list').empty();

    // Append the updated channel list to the HTML
    channelList.forEach(channel => {
        $('#channels-list').append(`<li class="channel" data-channel="${channel}">${channel}</li>`);
    });

    // Add click event listener to channel items
    $('.channel').click(function () {
        const channelName = $(this).data('channel');
        changeChatRoomText(channelName);
        clearChatMessages();
        startContinuousFetch(channelName);
        clearUserList(); // Clear user list when a new channel is clicked
        startContinuousUserFetch(channelName);
    });
}
    