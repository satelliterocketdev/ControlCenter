import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';
import MainMenu from '../../components/mainMenu/index';
import { setActiveMenuPosition } from '../../actions/menu/index';
import { logoutAction } from '../../actions/session/index';
import WorkerDashboard from "../../components/dashboard/worker";
import {isAllowedRole} from "../../utils/auth";
import {commonConstants} from "../../constants/common";
import Workers from "../../components/workers/list";
import {
    getMessages,
    getUsersUpdates,
    getFakeUsersUpdates,
    updateChatFakeUsers,
    selectChatFakeUser,
    changeFakeUserOnlineStatus,
    selectChat,
    changeChatMessage,
    sendChatMessage,
    readChatMessages,
    saveNotes,
    startNewChat
} from "../../actions/chat";
import { getWorkerFakeUsers } from "../../actions/fakeUsers";
import { getWorkersList, selectWorker, deleteWorker } from "../../actions/workers";
import {
    getUsersList,
    hideError,
    onSearchValueChange,
    searchUsersList,
    sortUsersList
} from "../../actions/users";

class DashboardView extends Component {
    componentDidMount() {

        if (this.props.menu.activeMenu !== 'dashboard') {
            this.props.setActiveMenuPosition('dashboard');
        }
    }

    render() {
        if (!this.props.session.isAuthenticated) {
            return <Redirect to={{pathname: '/signIn', state: {from: this.props.location}}}/>
        }
        return (
            <div>
                <MainMenu
                    needShowPrivateItems={this.props.session.isAuthenticated}
                    menu={ this.props.menu }
                    onLogout={ this.props.onLogout }
                    setActiveMenuPosition={ this.props.setActiveMenuPosition }
                    userDetails={this.props.session.userDetails}
                />
                {isAllowedRole([commonConstants.ADMIN_ROLE], this.props.session.userDetails) &&
                    <Workers
                        userDetails={this.props.session.userDetails}
                        workers={this.props.workers}
                        getWorkerFakeUsers={this.props.getWorkerFakeUsers}
                        getWorkersList={this.props.getWorkersList}
                        selectWorker={this.props.selectWorker}
                        deleteWorker={this.props.deleteWorker}
                    /> || <WorkerDashboard
                            session={this.props.session}
                            chat={this.props.chat}
                            users={this.props.users}
                            getMessages={this.props.getMessages}
                            getChatUsers={this.props.getChatUsers}
                            getFakeUsers={this.props.getFakeUsers}
                            getWorkerFakeUsers={this.props.getWorkerFakeUsers}
                            updateFakeUsers={this.props.updateFakeUsers}
                            selectChatFakeUser={this.props.selectChatFakeUser}
                            changeFakeUserOnlineStatus={this.props.changeFakeUserOnlineStatus}
                            selectChat={this.props.selectChat}
                            changeChatMessage={this.props.changeChatMessage}
                            sendChatMessage={this.props.sendChatMessage}
                            readChatMessages={this.props.readChatMessages}
                            saveNotes={this.props.saveNotes}
                            getUsersList={this.props.getUsersList}
                            hideError={this.props.hideError}
                            sortUsersList={this.props.sortUsersList}
                            onSearchValueChange={this.props.onSearchValueChange}
                            searchUsersList={this.props.searchUsersList}
                            startNewChat={this.props.startNewChat}
                        />
                }

            </div>
        );
    }
}

const mapStateToProps = (state) => {
    return {
        menu: state.menu,
        users: state.users,
        session: state.session,
        workers: state.workers,
        chat: state.chat,
    }
};

const mapDispatchToProps = (dispatch) => {
    return {
        setActiveMenuPosition(activeMenu) {
            setActiveMenuPosition(activeMenu)(dispatch)
        },
        onLogout() {
            logoutAction()(dispatch)
        },
        getMessages() {
            getMessages()(dispatch);
        },
        getChatUsers() {
            getUsersUpdates()(dispatch)
        },
        getFakeUsers() {
            getFakeUsersUpdates()(dispatch)
        },
        updateFakeUsers() {
            updateChatFakeUsers()(dispatch)
        },
        getWorkerFakeUsers(id) {
            getWorkerFakeUsers(id)(dispatch);
        },
        selectChatFakeUser(id) {
            selectChatFakeUser(id)(dispatch);
        },
        changeFakeUserOnlineStatus(uid, online) {
            changeFakeUserOnlineStatus(uid, online)(dispatch);
        },
        selectChat(chatId) {
            selectChat(chatId)(dispatch);
        },
        changeChatMessage(e, data) {
            changeChatMessage(data.value)(dispatch);
        },
        sendChatMessage(chatId, fromName, fromUserId, recipientToken, type, message) {
            sendChatMessage(chatId, fromName, fromUserId, recipientToken, type, message)(dispatch);
        },
        readChatMessages(uid, chatId, messages) {
            readChatMessages(uid, chatId, messages)(dispatch);
        },
        getWorkersList() {
            getWorkersList()(dispatch);
        },
        selectWorker(id) {
            selectWorker(id)(dispatch);
        },
        deleteWorker(id) {
            deleteWorker(id)(dispatch);
        },
        saveNotes(chatId, value) {
            saveNotes(chatId, value)(dispatch);
        },
        getUsersList: (offset, ordered, orderedField) => {
            getUsersList(offset, ordered, orderedField)(dispatch);
        },
        onSearchValueChange: (e, data) => {
            onSearchValueChange(data.value)(dispatch);
        },
        searchUsersList: (value) => {
            searchUsersList(value)(dispatch);
        },
        sortUsersList: (offset, prevColumn, column, prevDirection) => {
            sortUsersList(offset, prevColumn, column, prevDirection)(dispatch);
        },
        startNewChat: (fakeUserId, realUserId) => {
            startNewChat(fakeUserId, realUserId)(dispatch);
        },
        hideError: () => {
            hideError()(dispatch)
        },
    }
};

export default connect(mapStateToProps, mapDispatchToProps)(DashboardView);

