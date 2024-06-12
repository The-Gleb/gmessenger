document.addEventListener('DOMContentLoaded', function() {
    const chatList = document.querySelectorAll('.chat');

    chatList.forEach(chat => {
        chat.addEventListener('click', function() {

            location.replace("/static/");


            // const chatType = this.classList.contains('group-chat') ? 'group' : 'dialog';
            // const receiverID = this.dataset.receiverid;
            // const groupID = this.dataset.groupid;
            //
            // if (chatType === 'dialog') {
            //     fetch(`/dialog/${receiverID}`)
            //         .then(response => {
            //             if (response.ok) {
            //                 // Обработка успешного ответа
            //             } else {
            //                 throw new Error('Network response was not ok');
            //             }
            //         })
            //         .catch(error => {
            //             console.error('Error:', error);
            //         });
            // } else if (chatType === 'group') {
            //     fetch(`/group/${groupID}`)
            //         .then(response => {
            //             if (response.ok) {
            //                 // Обработка успешного ответа
            //             } else {
            //                 throw new Error('Network response was not ok');
            //             }
            //         })
            //         .catch(error => {
            //             console.error('Error:', error);
            //         });
            // }
        });
    });
});