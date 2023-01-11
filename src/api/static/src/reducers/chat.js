import { chatActionTypes } from '../constants/actions/chat';
import { fakeUsersActionTypes } from '../constants/actions/fakeUsers';

const defaultFakeChatUser = {
    id: 0,
    name: "",
    external_id: "",
    owner_id: 0,
    online: false
};


const initialState = {
    loading: false,
    error: {
        visible: false,
        header: '',
        messages: []
    },
    currentFakeUser: {
        id: 0,
        external_id: "",
    },
    allFakeUsers: {},
    currentFakeUsers: [],
    allChats: {},
    currentChat: null,
    chatMessage: ""
};

export default function chat(state = initialState, action) {
    switch (action.type) {
        case chatActionTypes.HIDE_ERROR:
            return {
                ...state,
                error: {
                    ...state.error,
                    visible: false
                }
            };

        case fakeUsersActionTypes.GET_WORKER_FAKE_USERS_REQUEST:
            return {
                ...state,
                loading: true
            };

        case fakeUsersActionTypes.GET_WORKER_FAKE_USERS_FAILURE:
            return {
                ...state,
                loading: false
            };

        case fakeUsersActionTypes.GET_WORKER_FAKE_USERS_SUCCESS:
            return {
                ...state,
                currentFakeUsers: action.data.fake_users.map((item, index) => {return {...defaultFakeChatUser, ...item}}),
                loading: false
            };

        case fakeUsersActionTypes.GET_ALL_FAKE_USERS_UPDATE:
            return {
                ...state,
                allFakeUsers: action.data
            };

        case chatActionTypes.CHAT_NEED_UPDATE_FAKE_USERS:
            return {
                ...state,
                currentFakeUsers: state.currentFakeUsers.map((item, index) => {
                    if (!state.allFakeUsers[item.external_id]) {
                        return item;
                    }

                    const fakeUser = state.allFakeUsers[item.external_id];
                    return {
                        ...item,
                        online: fakeUser.online || false
                    };
                })
            };

        case chatActionTypes.CHAT_SELECT_FAKE_USER:
            const currentFakeUser = state.currentFakeUsers.find((item) => item.id === action.data);
            if (!currentFakeUser) {
                return state;
            }
            return {
                ...state,
                currentFakeUser: {
                    ...state.currentFakeUser,
                    id: currentFakeUser.id,
                    external_id: currentFakeUser.external_id
                },
                currentChat: null
            };

        case chatActionTypes.CHAT_CHANGE_FAKE_USER_STATUS_REQUEST:
            return {
                ...state,
                currentFakeUsers: state.currentFakeUsers.map((item, index) => {
                    if (action.data.uid !== item.external_id) {
                        return item;
                    }

                    return {
                        ...item,
                        online: action.data.online
                    }
                })
            };

        case chatActionTypes.CHAT_CHANGE_FAKE_USER_STATUS_FAILURE:
            return {
                ...state,
                currentFakeUsers: state.currentFakeUsers.map((item, index) => {
                    if (action.data.uid !== item.external_id) {
                        return item;
                    }

                    return {
                        ...item,
                        online: !action.data.online
                    }
                })
            };


        case chatActionTypes.SAVE_CHAT_NOTES_FAILURE:
            return {
                ...state,
                error: {
                    ...state.error,
                    header: action.data.header && action.data.header !== '' ? action.data.header : '',
                    messages: action.data.errors && action.data.errors.length > 0 ? action.data.errors : [],
                    visible: true
                },
                loading: false
            };

        case chatActionTypes.CHAT_CHANGE_FAKE_USER_STATUS_SUCCESS:
            return {
                ...state,
                currentFakeUsers: state.currentFakeUsers.map((item, index) => {
                    if (action.data.uid !== item.external_id) {
                        return item;
                    }

                    return {
                        ...item
                    }
                })
            };

        /*case chatActionTypes.READ_CHAT_MESSAGES_REQUEST:
            return {
                ...state,
                allChats: {
                    ...state.allChats,
                    [action.data.chatId]: {
                        ...state.allChats[action.data.chatId],
                        history: {
                            ...Object.entries(state.allChats[action.data.chatId].history).reduce((acc, [key, value]) => {
                                if (action.data.messages.some((item) => item === key)) {
                                    return {
                                        ...value,
                                        ir: true
                                    }
                                }
                                return value;
                            }, {})
                        }
                    }
                }
            };*/

        case chatActionTypes.GET_CHAT_MESSAGES_UPDATE:
            if (action.data.path === "/") {
                return {
                    ...state,
                    allChats: action.data.data
                };
            }

            var pathParts = action.data.path.split('/');
            if (pathParts.length < 2) {
                return state;
            }

            var existingChat  = state.allChats[pathParts[1]];
            if (!existingChat) {
                return {
                    ...state,
                    allChats: {
                        ...state.allChats,
                        [pathParts[1]]: {
                            history: {}
                        }
                    }
                };
            }

            if (action.data.path.includes('notes')) {
                return {
                    ...state,
                    allChats: {
                        ...state.allChats,
                        [pathParts[1]]: {
                            ...state.allChats[pathParts[1]],
                            notes: action.data.data
                        }
                    }
                };
            }

            var existingChat = state.allChats[pathParts[1]];
            if (!existingChat) {
                return {
                    ...state,
                    allChats: {
                        ...state.allChats,
                        [pathParts[1]]: {
                            ...state.allChats[pathParts[1]],
                            history: {
                                [pathParts[3]]: {
                                    ...action.data.data
                                }
                            }
                        }
                    }
                };
            }
            if (pathParts.length === 4) {
                return {
                    ...state,
                    allChats: {
                        ...state.allChats,
                        [pathParts[1]]: {
                            ...state.allChats[pathParts[1]],
                            history: {
                                ...state.allChats[pathParts[1]].history,
                                [pathParts[3]]: {
                                    ...state.allChats[pathParts[1]].history[pathParts[3]],
                                    ...action.data.data
                                }
                            }
                        }
                    }
                };
            }
            if (pathParts.length === 5 && pathParts[4] === "ir") {
                return {
                    ...state,
                    allChats: {
                        ...state.allChats,
                        [pathParts[1]]: {
                            ...state.allChats[pathParts[1]],
                            history: {
                                ...state.allChats[pathParts[1]].history,
                                [pathParts[3]]: {
                                    ...state.allChats[pathParts[1]].history[pathParts[3]],
                                    ir: action.data.data
                                }
                            }
                        }
                    }
                };
            }
            return state;

        case chatActionTypes.GET_CHAT_MESSAGES_NEW:
            if (action.data.path !== "/") {
                return state;
            }

            const key = Object.keys(action.data.data)[0];

            var existingChat  = state.allChats[key];
            if (!existingChat) {
                return {
                    ...state,
                    allChats: {
                        ...state.allChats,
                        [key]: {
                            history: {}
                        }
                    }
                };
            }


            return state;

        case chatActionTypes.GET_CHAT_USERS_UPDATE:
            if (action.data.path === "/") {
                return {
                    ...state,
                    allUsers: action.data.data
                };
            }

            if (!action.data.path.includes('online')) {
                return state;
            }

            var pathParts = action.data.path.split('/');
            return {
                ...state,
                allUsers: {
                    ...state.allUsers,
                    [pathParts[1]]: {
                        ...state.allUsers[pathParts[1]],
                        online: action.data.data
                    }
                }
            };

        case chatActionTypes.SELECT_CHAT:
            return {
                ...state,
                currentChat: action.data
            };
        case chatActionTypes.CHANGE_CHAT_MESSAGE:
            return {
                ...state,
                chatMessage: action.data
            };

        case chatActionTypes.SEND_CHAT_MESSAGE_REQUEST:
            return {
                ...state,
                chatMessage: "",
                loading: false
            };

        case chatActionTypes.SEND_CHAT_MESSAGE_FAILURE:
            return {
                ...state,
                error: {
                    ...state.error,
                    header: action.data.header && action.data.header !== '' ? action.data.header : '',
                    messages: action.data.errors && action.data.errors.length > 0 ? action.data.errors : [],
                    visible: true
                },
                loading: false
            };

        case chatActionTypes.SEND_CHAT_MESSAGE_SUCCESS:
            return {
                ...state,
                loading: false
            };

        default:
            return state;
    }

}
