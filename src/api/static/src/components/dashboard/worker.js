import React, { Component } from 'react';
import {Segment, Icon, Header, Dimmer, Loader, Label, List, Grid, Button, Form, Divider, Popup, Checkbox} from 'semantic-ui-react';

import ErrorMessage from '../errorMessage/index';
import PropTypes from "prop-types";
import {Link} from "react-router-dom";
import AddNewPhoto from "../imageUpload/addNewPhoto";
import Photo from "../imageUpload/photo";
import Notes from "../chat/notes";
import ChatUsers from "../chat/userList";


export default class WorkerDashboard extends Component {
    messagesEnd = null;
    state = {usersList: false};
    scrollToBottom = () => {
        if (this.messagesEnd) {
            this.messagesEnd.scrollIntoView({ behavior: "smooth" });
        }
    };

    onStartNewChat = (fakeUserId, realUserId) => {
        this.props.startNewChat(fakeUserId, realUserId);
        this.closeUsersList();
    }

    openUsersList() {
        this.setState({usersList: true});
    }

    closeUsersList() {
        this.setState({usersList: false});
    }



    componentDidMount() {
        this.props.getMessages();
        this.props.getChatUsers();
        this.props.getWorkerFakeUsers(this.props.session.userDetails.id);
        this.scrollToBottom();
    }

    componentDidUpdate(prevProps) {
        if(this.props.chat.currentFakeUsers.length > 0 && prevProps.chat.currentFakeUsers.length === 0) {
            this.props.getFakeUsers();
        }

        if (this.props.chat.allFakeUsers !== prevProps.chat.allFakeUsers) {
            this.props.updateFakeUsers();
        }
        this.scrollToBottom();
    }

    componentWillReceiveProps(prevProps) {
        if (this.props.chat.allFakeUsers !== prevProps.chat.allFakeUsers) {
            this.props.updateFakeUsers();
        }
    }

    render() {
        const chats = Object.entries(this.props.chat.allChats)
            .filter(([key, value]) => this.props.chat.currentFakeUser.external_id && key.includes(this.props.chat.currentFakeUser.external_id) && key.split('___').length === 2)
            .map(([key, value]) => {
                const chatParts = key.split('___');
                const userId = chatParts[0] !== this.props.chat.currentFakeUser.external_id ? chatParts[0] : chatParts[1];
                return {
                    fakeUserId: this.props.chat.currentFakeUser.external_id,
                    chatId: key,
                    userId: userId,
                    history: value.history,
                    notes: value.notes
                };
            });

        const allPreparedChats = Object.entries(this.props.chat.allChats)
            .filter(([key, value]) => key.split('___').length === 2 && key.includes('x0x'))
            .map(([key, value]) => {
                const chatParts = key.split('___');
                const userId = !chatParts[0].includes('x0x') ? chatParts[0] : chatParts[1];
                const fakeUserId = chatParts[0].includes('x0x') ? chatParts[0] : chatParts[1];
                return {
                    fakeUserId: fakeUserId,
                    chatId: key,
                    userId: userId,
                    history: value.history
                };
            });
        const currentChat = chats.find((item, index) => item.chatId === this.props.chat.currentChat);
        if (currentChat) {
            let unreadMessages = currentChat.history
                ? Object.entries(currentChat.history).filter(([key, value]) => !value.ir && value.uid !== this.props.chat.currentFakeUser.external_id).map(([key, value]) => key)
                : [];

            if (unreadMessages && unreadMessages.length > 0) {
                this.props.readChatMessages(this.props.chat.currentFakeUser.external_id, currentChat.chatId, [unreadMessages[0]]);
            }
        }

        const currentFakeUser = this.props.chat.currentFakeUsers.find((item, index) => item.id === this.props.chat.currentFakeUser.id);
        const currentRealUser = currentChat && this.props.chat.allUsers[currentChat.userId];

        return (
            <div>
                <Dimmer.Dimmable dimmed={this.props.chat.loading}>
                    <Dimmer active={this.props.chat.loading} inverted>
                        <Loader>Loading</Loader>
                    </Dimmer>
                    {
                        this.props.chat.error && this.props.chat.error.visible &&
                        <ErrorMessage error={this.props.chat.error} hideError={this.props.hideError} />
                    }
                    <Grid columns={2}>
                        <Grid.Row stretched>
                            <Grid.Column width={6} floated='right'>
                                <Segment>
                                    <Grid columns={2}>
                                        <Grid.Column width={8}>
                                            <Header content='Moderators'/>
                                            <List divided relaxed selection verticalAlign='middle'>
                                                {this.props.chat.currentFakeUsers.map((item, index) => {
                                                    const currentUnreadMessages = allPreparedChats
                                                        .filter((c, index) => c.fakeUserId === item.external_id)
                                                        .reduce((acc, uc) => {
                                                            if (!uc.history) {
                                                                return acc;
                                                            }
                                                            return acc + Object.entries(uc.history).filter(([key, value]) => !value.ir && value.uid !== item.external_id).length;
                                                        }, 0);
                                                    const isActive = this.props.chat.currentFakeUser.external_id === item.external_id;
                                                    return <List.Item
                                                        key={index}
                                                        active={isActive}
                                                        onClick={(e) => this.props.selectChatFakeUser(item.id)}>
                                                        {<List.Content floated='right'>
                                                            <Link to={"/workers/"+this.props.session.userDetails.id+"/fakeUsers/edit/"+item.external_id} >
                                                                <Icon link name="pencil" size="large" color="blue"/>
                                                            </Link>
                                                        </List.Content> }
                                                        {<List.Content floated='right'>
                                                            <Popup
                                                                trigger={<Icon
                                                                    link
                                                                    name="chat"
                                                                    size="large"
                                                                    color={(this.state.usersList && isActive) ? 'yellow': 'blue'}
                                                                    onClick={(e) => {
                                                                        (this.state.usersList && isActive) ? this.closeUsersList() : this.openUsersList()
                                                                    }}/>}
                                                                content='Start new chat'
                                                                position='left center'
                                                            />

                                                        </List.Content>}
                                                        <List.Content floated='right'>
                                                            <Checkbox
                                                                slider
                                                                checked={item.online}
                                                                onChange={(e) => {e.stopPropagation(); this.props.changeFakeUserOnlineStatus(item.external_id, !item.online)}}
                                                            />
                                                        </List.Content>
                                                        <List.Content floated='right'>
                                                            <Label color={currentUnreadMessages > 0 ? 'yellow' : 'grey'}>{currentUnreadMessages}</Label>
                                                        </List.Content>
                                                        <List.Content>
                                                            <List.Header>
                                                                {item.name}
                                                            </List.Header>
                                                        </List.Content>
                                                    </List.Item>
                                                })}
                                            </List>
                                            {false && <Grid>
                                                <Grid.Column textAlign="center">
                                                    <Link to={"/workers/"+this.props.session.userDetails.id+"/fakeUsers/new"}>
                                                        <Button circular icon='add' positive/>
                                                    </Link>
                                                </Grid.Column>
                                            </Grid>}
                                        </Grid.Column>
                                        <Grid.Column width={8}>
                                            <Header content='Users'/>
                                            <List divided relaxed selection>
                                                {chats.map((item, index) => {
                                                    const user = this.props.chat.allUsers[item.userId];
                                                    const currentUnreadMessages = allPreparedChats
                                                        .filter((c, index) => c.chatId === item.chatId)
                                                        .reduce((acc, uc) => {
                                                            if (!uc.history) {
                                                                return acc;
                                                            }
                                                            return acc + Object.entries(uc.history).filter(([key, value]) => !value.ir && value.uid !== item.fakeUserId).length;
                                                        }, 0);
                                                    return <List.Item key={index} onClick={(e) => {
                                                        this.props.selectChat(item.chatId);
                                                    }}
                                                    active={this.props.chat.currentChat && this.props.chat.currentChat === item.chatId}>
                                                        {<List.Content floated='right'>
                                                            <Link target='_blank' to={"/users/edit/"+item.userId}>
                                                                <Icon link name="pencil" size="large" color="blue"/>
                                                            </Link>
                                                        </List.Content>}
                                                        <List.Content floated='right'>
                                                            <Icon
                                                                name='circle'
                                                                size='tiny'
                                                                color={(user && user.online) ? 'green' : 'yellow'}
                                                                corner
                                                            />
                                                        </List.Content>
                                                        <List.Content floated='right'>
                                                            <Label color={currentUnreadMessages > 0 ? 'yellow' : 'grey'}>{currentUnreadMessages}</Label>
                                                        </List.Content>
                                                        <List.Content>
                                                            <List.Header>{(user && user.name) ? user.name : "Undefined"}</List.Header>
                                                            {user && user.filters && user.filters.location && user.filters.location.lat && user.filters.location.lon &&
                                                            <List.Description>{"Lat: "+user.filters.location.lat}</List.Description>}
                                                            {user && user.filters && user.filters.location && user.filters.location.lat && user.filters.location.lon &&
                                                            <List.Description>{"Lon: " + user.filters.location.lon }</List.Description>}
                                                        </List.Content>

                                                    </List.Item>
                                                })}
                                            </List>
                                        </Grid.Column>
                                    </Grid>
                                </Segment>
                            </Grid.Column>
                            <Grid.Column width={10} >
                                {this.state.usersList && currentFakeUser && (
                                    <div>
                                        <Button icon='close' onClick={e => this.closeUsersList()} />
                                        <ChatUsers
                                            fakeUserId={this.props.chat.currentFakeUser.external_id}
                                            startNewChat={this.onStartNewChat}
                                            chat={this.props.chat}
                                            users={ this.props.users }
                                            getUsersList={this.props.getUsersList}
                                            hideError={this.props.hideError}
                                            sortUsersList={this.props.sortUsersList}
                                            onSearchValueChange={this.props.onSearchValueChange}
                                            searchUsersList={this.props.searchUsersList}/>
                                    </div>)}
                                {!this.state.usersList && <Segment style={{maxHeight: '500px', height: '500px', overflowX: 'scroll'}}  color='teal'>

                                    <Grid columns={2}>
                                        {currentChat && currentChat.history && Object.entries(currentChat.history).map(([key, value]) => {
                                            const date = new Date(value.ts);
                                            if (value.uid === currentChat.userId) {
                                                const user = this.props.chat.allUsers[currentChat.userId];
                                                return <Grid.Row textAlign={value.url ? 'center' : 'left'} key={key}>
                                                    <Grid.Column>
                                                        <Popup
                                                            trigger={<Segment tertiary size='tiny'>
                                                                {value.url &&
                                                                <a target='_blank' href={value.url}>
                                                                    <Photo src={value.url} height={'100px'}/>
                                                                </a>}
                                                                {value.txt && value.txt}

                                                                <Divider hidden fitted/>
                                                                {value.ir && <Header icon={<Icon corner size='mini' name='check'/>} disabled sub floated='right' content={date.toLocaleString()}/> || <Header disabled sub floated='right' content={date.toLocaleString()}/>}
                                                            </Segment>}
                                                            content={(user && user.name) ? user.name : 'Undefined'}
                                                            basic
                                                            mouseEnterDelay={1000}
                                                            on='hover'
                                                            size={'mini'}
                                                            position='bottom right'
                                                        />

                                                    </Grid.Column>
                                                    <Grid.Column>
                                                    </Grid.Column>
                                                </Grid.Row>
                                            } else {
                                                const user = this.props.chat.allFakeUsers[value.uid];
                                                return <Grid.Row textAlign={value.url ? 'center' : 'left'} key={key}>
                                                    <Grid.Column>
                                                    </Grid.Column>
                                                    <Grid.Column>
                                                        <Popup
                                                            trigger={<Segment inverted color={'teal'} size='tiny'>
                                                                {value.url &&
                                                                <a target='_blank' href={value.url}>
                                                                    <Photo src={value.url} height={'100px'}/>
                                                                </a>}
                                                                {value.txt && value.txt}

                                                                <Divider hidden fitted/>
                                                                {value.ir && <Header icon={<Icon corner size='mini' name='check'/>} sub floated='right' content={date.toLocaleString()}/> || <Header sub floated='right' content={date.toLocaleString()}/>}
                                                            </Segment>}
                                                            content={(user && user.name) ? user.name : 'Undefined'}
                                                            basic
                                                            mouseEnterDelay={1000}
                                                            on='hover'
                                                            size={'mini'}
                                                            position='bottom right'
                                                        />
                                                    </Grid.Column>
                                                </Grid.Row>
                                            }

                                        })}
                                    </Grid>
                                    <div style={{ float:"left", clear: "both" }}
                                         ref={(el) => { this.messagesEnd = el; }}/>
                                </Segment>}
                                {!this.state.usersList && currentChat && <Notes id={currentChat.chatId} value={currentChat.notes || ''} saveValue={this.props.saveNotes}/>}
                                <Divider hidden/>
                                {!this.state.usersList && <Grid columns={2}>
                                    <Grid.Column>
                                        <Form reply>
                                            <Form.TextArea rows={5} value={this.props.chat.chatMessage} onChange={this.props.changeChatMessage} disabled={!this.props.chat.currentChat}/>
                                            <Divider hidden/>
                                            <Grid textAlign='center'>

                                                <Button
                                                    content='Send'
                                                    primary
                                                    //chatId, fromName, fromUserId, recipientToken, message
                                                    onClick={(e) => {e.preventDefault();
                                                        this.props.sendChatMessage(
                                                            this.props.chat.currentChat,
                                                            currentFakeUser.name,
                                                            currentFakeUser.external_id,
                                                            currentRealUser.fcmToken,
                                                            'text',
                                                            this.props.chat.chatMessage
                                                        )}}
                                                    disabled={!this.props.chat.currentChat || this.props.chat.chatMessage === ""}/>
                                            </Grid>
                                        </Form>
                                    </Grid.Column>
                                    <Grid.Column>
                                        <AddNewPhoto folder={currentFakeUser ? currentFakeUser.external_id : 'images'} disabled={!this.props.chat.currentChat} height={'100px'} header={'Send image'} onAdd={(url) => {
                                            this.props.sendChatMessage(
                                                this.props.chat.currentChat,
                                                currentFakeUser.name,
                                                currentFakeUser.external_id,
                                                currentRealUser.fcmToken,
                                                'image',
                                                url
                                            )}}/>
                                    </Grid.Column>
                                </Grid>}

                            </Grid.Column>
                        </Grid.Row>
                    </Grid>

                </Dimmer.Dimmable>
            </div>
        )
    }
}

WorkerDashboard.propTypes = {
    chat: PropTypes.shape({
        loading: PropTypes.bool,
        error:  PropTypes.shape({
            visible: PropTypes.bool
        })
    }).isRequired,
    session: PropTypes.shape({
        userDetails: PropTypes.shape({
            id: PropTypes.number.isRequired
        }).isRequired
    }).isRequired
};
