import { chatActionTypes } from "../../constants/actions/chat";
import { Api } from "../../utils/api";
import {fakeUsersActionTypes} from "../../constants/actions/fakeUsers";
import {sessionActionTypes} from "../../constants/actions/session";
import {usersActionTypes} from "../../constants/actions/users";

const messageHeader = 'API request failed';

function onSSEFailure(action, dispatch) {
    return function (error) {
        dispatch({
            type: action,
            data: error
        });
    }
}

function onSSEData(action, dispatch) {
    return function (data) {
        dispatch({
            type: action,
            data: data
        });
    }
}

export function getMessages() {
    return function(dispatch) {
        dispatch({
            type: chatActionTypes.GET_CHAT_MESSAGES_LIST_REQUEST
        });
        Api.getChatMessages(onSSEData(chatActionTypes.GET_CHAT_MESSAGES_UPDATE, dispatch), onSSEData(chatActionTypes.GET_CHAT_MESSAGES_NEW, dispatch), onSSEFailure(chatActionTypes.GET_CHAT_MESSAGES_LIST_FAILURE, dispatch));
    }
}

export function getUsersUpdates() {
    return function (dispatch) {
        dispatch({
            type: chatActionTypes.GET_CHAT_USERS_LIST_REQUEST
        });
        Api.getChatUsersDetails(onSSEData(chatActionTypes.GET_CHAT_USERS_UPDATE, dispatch), onSSEFailure(chatActionTypes.GET_CHAT_USERS_LIST_FAILURE, dispatch));
    }
}

export function getFakeUsersUpdates() {
    return function(dispatch) {
        dispatch({
            type: fakeUsersActionTypes.GET_ALL_FAKE_USERS_LIST_REQUEST
        });
        return Api.getFakeUsersDetails().then(response => {
            const errorMessage = "Unknown error. HTTP Status code: " + response.status;
            if (response.status !== 200) {
                dispatch({
                    type: chatActionTypes.GET_ALL_FAKE_USERS_LIST_FAILURE,
                    data: {errors: [errorMessage], header: messageHeader}
                });
                return;
            }


            response.json().then(json => {
                dispatch({
                    type: fakeUsersActionTypes.GET_ALL_FAKE_USERS_UPDATE,
                    data: json
                });
            }).catch(error => {
                console.log('fail', error);
                dispatch({
                    type: fakeUsersActionTypes.GET_ALL_FAKE_USERS_LIST_FAILURE,
                    data: {header: messageHeader, errors: [error]}
                });
            });
        }).catch(error => {
            console.log('fail', error);
            dispatch({
                type: fakeUsersActionTypes.GET_ALL_FAKE_USERS_LIST_FAILURE,
                data: {header: messageHeader, errors: [error]}
            });
        });
    };
}

export function updateChatFakeUsers() {
    return function (dispatch) {
        dispatch({
            type: chatActionTypes.CHAT_NEED_UPDATE_FAKE_USERS
        });
    }
}

export function selectChatFakeUser(id) {
    return function (dispatch) {
        dispatch({
            type: chatActionTypes.CHAT_SELECT_FAKE_USER,
            data: id
        });
    }
}

export function selectChat(chatId) {
    return function (dispatch) {
        dispatch({
            type: chatActionTypes.SELECT_CHAT,
            data: chatId
        });
    }
}

export function changeFakeUserOnlineStatus(uid, online) {
    return function(dispatch) {
        dispatch({
            type: chatActionTypes.CHAT_CHANGE_FAKE_USER_STATUS_REQUEST,
            data: {uid: uid, online: online}
        });
        return Api.changeFakeUserOnlineStatus(uid, online).then(response => {
            if (response.status === 401) {
                dispatch({
                    type: chatActionTypes.CHAT_CHANGE_FAKE_USER_STATUS_FAILURE,
                    data: {uid: uid, online: online}
                });
                dispatch({
                    type: sessionActionTypes.LOGOUT_REQUEST,
                });
                return;
            }

            if (response.status === 404) {
                dispatch({
                    type:chatActionTypes.CHAT_CHANGE_FAKE_USER_STATUS_FAILURE,
                    data: {uid: uid, online: online}
                });
                return;
            }

            if (response.status !== 200) {
                dispatch({
                    type: chatActionTypes.CHAT_CHANGE_FAKE_USER_STATUS_FAILURE,
                    data: {uid: uid, online: online}
                });
                return;
            }

            dispatch({
                type: chatActionTypes.CHAT_CHANGE_FAKE_USER_STATUS_SUCCESS,
                data: {uid: uid, online: online}
            });


        }).catch(error => {
            console.log('fail', error);
            dispatch({
                type: usersActionTypes.EDIT_USER_DETAILS_FAILURE,
                data: {uid: uid, online: online}
            });
        });
    }
}

export function changeChatMessage(chatMessage) {
    return function (dispatch) {
        dispatch({
            type: chatActionTypes.CHANGE_CHAT_MESSAGE,
            data: chatMessage
        });
    }
}

export function sendChatMessage(chatId, fromName, fromUserId, recipientToken, type, message) {
    return function(dispatch) {
        dispatch({
            type: chatActionTypes.SEND_CHAT_MESSAGE_REQUEST
        });
        return Api.sendChatMessage(chatId, fromName, fromUserId, recipientToken, type, message).then(response => {
            const errorMessage = "Unknown error. HTTP Status code: " + response.status;
            if (response.status === 401) {
                dispatch({
                    type: chatActionTypes.SEND_CHAT_MESSAGE_FAILURE,
                    data: {header: messageHeader}
                });
                dispatch({
                    type: sessionActionTypes.LOGOUT_REQUEST,
                });
                return;
            }

            if (response.status !== 200) {
                dispatch({
                    type: chatActionTypes.SEND_CHAT_MESSAGE_FAILURE,
                    data: {errors: [errorMessage], header: messageHeader}
                });
                return;
            }

            dispatch({
                type: chatActionTypes.SEND_CHAT_MESSAGE_SUCCESS
            });


        }).catch(error => {
            console.log('fail', error);
            dispatch({
                type: chatActionTypes.SEND_CHAT_MESSAGE_FAILURE,
                data: {header: messageHeader}
            });
        });
    }
}


export function readChatMessages(uid, chatId, messages) {
    return function(dispatch) {
        dispatch({
            type: chatActionTypes.READ_CHAT_MESSAGES_REQUEST,
            data: {messages, chatId}
        });
        return Api.readChatMessages(uid, chatId, messages).then(response => {
            const errorMessage = "Unknown error. HTTP Status code: " + response.status;
            if (response.status === 401) {
                dispatch({
                    type: chatActionTypes.READ_CHAT_MESSAGES_FAILURE,
                    data: {header: messageHeader}
                });
                dispatch({
                    type: sessionActionTypes.LOGOUT_REQUEST,
                });
                return;
            }

            if (response.status !== 200) {
                dispatch({
                    type: chatActionTypes.READ_CHAT_MESSAGES_FAILURE,
                    data: {errors: [errorMessage], header: messageHeader}
                });
                return;
            }

            dispatch({
                type: chatActionTypes.READ_CHAT_MESSAGES_SUCCESS
            });


        }).catch(error => {
            console.log('fail', error);
            dispatch({
                type: chatActionTypes.READ_CHAT_MESSAGES_FAILURE,
                data: {header: messageHeader}
            });
        });
    }
}

export function saveNotes(chatId, notes) {
    return function(dispatch) {
        return Api.saveChatNotes(chatId, notes).then(response => {
            const errorMessage = "Unknown error. HTTP Status code: " + response.status;
            if (response.status === 401) {
                dispatch({
                    type: sessionActionTypes.LOGOUT_REQUEST,
                });
                return;
            }

            if (response.status !== 200) {
                dispatch({
                    type: chatActionTypes.SAVE_CHAT_NOTES_FAILURE,
                    data: {errors: [errorMessage], header: messageHeader}
                });
            }


        }).catch(error => {
            console.log('fail', error);
            dispatch({
                type: chatActionTypes.SAVE_CHAT_NOTES_FAILURE,
                data: {header: messageHeader}
            });
        });
    }
}


export function startNewChat(fakeUserId, realUserId) {
    return function(dispatch) {
        return Api.startNewChat(fakeUserId, realUserId).then(response => {
            const errorMessage = "Unknown error. HTTP Status code: " + response.status;
            if (response.status === 401) {
                dispatch({
                    type: sessionActionTypes.LOGOUT_REQUEST,
                });
                return;
            }

            if (response.status !== 200) {
                dispatch({
                    type: chatActionTypes.START_NEW_CHAT_FAILURE,
                    data: {errors: [errorMessage], header: messageHeader}
                });
            }


        }).catch(error => {
            console.log('fail', error);
            dispatch({
                type: chatActionTypes.START_NEW_CHAT_FAILURE,
                data: {header: messageHeader}
            });
        });
    }
}
