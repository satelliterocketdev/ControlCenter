import React, { Component } from 'react';
import { Link, Redirect } from 'react-router-dom';
import { Form, Segment, Button, Dropdown, Dimmer, Loader, Image, Header, Icon, Card } from 'semantic-ui-react';
import ErrorMessage from '../errorMessage/index';
import {commonConstants} from "../../constants/common";

const genderOptions = [
    { key: 'm', text: 'Male', value: 'male' },
    { key: 'f', text: 'Female', value: 'female' }
];

const defaultDetails = [
    "name", "age", "gender", "about", "photos"
];

const detailsBlackList = [
    "filters", "matches", "likes"
];

export default class User extends Component {

    render() {
        if (this.props.users.currentUser._updated) {
            return <Redirect to={{pathname: '/users'}}/>
        }
        return (
        <Segment.Group>
            <Dimmer.Dimmable as={Segment} dimmed={this.props.users.loading}>
                <Dimmer active={this.props.users.loading} inverted>
                    <Loader>Loading</Loader>
                </Dimmer>
                {
                    this.props.users.error && this.props.users.error.visible &&
                    <ErrorMessage error={this.props.users.error} hideError={this.props.hideError} />
                }
                <Segment secondary textAlign='center'>
                    <Button onClick={(e, data) => this.props.unblockUser(this.props.users.currentUser.main.uid)}
                            size='large'
                            icon='unlock'
                            color='green'
                            disabled={(!!this.props.users.currentUser.main && !this.props.users.currentUser.main.disabled)}
                            content='Unblock user'/>
                    <Button onClick={(e, data) => this.props.blockUser(this.props.users.currentUser.main.uid)}
                            size='large'
                            icon='lock'
                            color='yellow'
                            disabled={(!!this.props.users.currentUser.main && this.props.users.currentUser.main.disabled)}
                            content='Block user'/>
                    <Button onClick={(e, data) => this.props.removeUserReports(this.props.users.currentUser.main.uid)}
                            size='large'
                            icon='trash'
                            color='green'
                            disabled={(!!this.props.users.currentUser.main && this.props.users.currentUser.main.reports === 0)}
                            content='Remove reports'/>
                </Segment>
            <Segment padded>
                <Form>
                <Segment secondary>

                        <Form.Input
                            fluid
                            label='Name'
                            placeholder='Name'
                            value={this.props.users.currentUser.main && this.props.users.currentUser.main.display_name ? this.props.users.currentUser.main.display_name : ''}
                            onChange={this.props.changeName}
                            />
                        <Form.Input
                            fluid
                            label='Email'
                            placeholder='Email'
                            value={this.props.users.currentUser.main && this.props.users.currentUser.main.email ? this.props.users.currentUser.main.email : ''}
                            onChange={this.props.changeEmail}
                            />
                        <Form.Input
                            fluid
                            label='Phone'
                            placeholder='Phone'
                            value={this.props.users.currentUser.main && this.props.users.currentUser.main.phone ? this.props.users.currentUser.main.phone : ''}
                            onChange={this.props.changePhone}
                            />
                        <Form.Input
                            fluid
                            label='Disabled'
                            placeholder='Disabled'
                            disabled={true}
                            readOnly={true}
                            value={this.props.users.currentUser.main && this.props.users.currentUser.main.disabled ? 'true' : 'false'}
                            />
                        <Form.Input
                            fluid
                            label='Reports'
                            placeholder='Reports'
                            disabled={true}
                            readOnly={true}
                            value={this.props.users.currentUser.main && this.props.users.currentUser.main.reports ? this.props.users.currentUser.main.reports : 0}
                        />
                        <Form.Input
                            fluid
                            label='SignUp date'
                            placeholder='SignUp date'
                            disabled={true}
                            readOnly={true}
                            value={this.props.users.currentUser.main && this.props.users.currentUser.main.sign_up ? (new Date(this.props.users.currentUser.main.sign_up)).toLocaleString() : '-'}
                        />
                        {this.props.users.currentUser.details && this.props.users.currentUser.details.filters && this.props.users.currentUser.details.filters.location &&
                        <Form.Group widths='equal'>
                            <Form.Input
                                fluid
                                label='Longitude'
                                placeholder='Longitude'
                                disabled={true}
                                readOnly={true}
                                value={this.props.users.currentUser.details.filters.location.lon || ''}
                            />
                            <Form.Input
                                fluid
                                label='Latitude'
                                placeholder='Latitude'
                                disabled={true}
                                readOnly={true}
                                value={this.props.users.currentUser.details.filters.location.lat || ''}
                            />
                        </Form.Group>}
                        <Form.Input
                            fluid
                            label='Age'
                            placeholder='Age'
                            type='number'
                            min={18}
                            max={59}
                            value={this.props.users.currentUser.details && this.props.users.currentUser.details.age ? this.props.users.currentUser.details.age : 18}
                            onChange={this.props.changeAge}
                            />
                        <Form.Select
                            fluid
                            label='Gender'
                            options={genderOptions}
                            placeholder='Gender'
                            onChange={this.props.changeGender}
                            value={this.props.users.currentUser.details && this.props.users.currentUser.details.gender ? this.props.users.currentUser.details.gender : ''}/>

                        <Form.TextArea
                            label='About'
                            placeholder='About'
                            value={this.props.users.currentUser.details && this.props.users.currentUser.details.about ? this.props.users.currentUser.details.about : ''}
                            onChange={this.props.changeAbout}
                            />
                        <Header size='large'>Additional fields</Header>
                        {this.props.users.currentUser && this.props.users.currentUser.details && Object.keys(this.props.users.currentUser.details).map((key, index)=>{
                            if (defaultDetails.indexOf(key) !== -1) {
                                return;
                            }
                            if (detailsBlackList.indexOf(key) !== -1) {
                                return;
                            }
                            var value = this.props.users.currentUser.details[key];
                            var type = typeof this.props.users.currentUser.details[key];
                            switch (type) {
                                case "string":
                                    return <Form.Input
                                        key={index}
                                        fluid
                                        label={key}
                                        placeholder={key}
                                        readOnly={true}
                                        disabled
                                        value={value}
                                        />
                                case "number":
                                    return <Form.Input
                                        key={index}
                                        fluid
                                        type='number'
                                        label={key}
                                        placeholder={key}
                                        readOnly={true}
                                        disabled
                                        value={value}
                                        />
                                case "object":
                                    var options = [];
                                    var values = [];
                                    Object.keys(value).map((property, index) => {
                                        options.push({key: property, text: value[property], value: value[property]});
                                        values.push(value[property]);
                                    });

                                    if (key === 'interestedIn') {
                                        options = commonConstants.INTERESTED_IN_OPTIONS;
                                    }
                                    return <Form.Select
                                        key={index}
                                        fluid
                                        label={key}
                                        placeholder={key}
                                        readOnly={true}
                                        disabled
                                        multiple selection
                                        options={options}
                                        value={values}
                                        />
                            }
                        })}

                </Segment>
                <Segment secondary>
                    <Header size='large'>Photos</Header>
                    <Card.Group>
                        {this.props.users.currentUser.details &&
                        this.props.users.currentUser.details.photos &&
                        this.props.users.currentUser.details.photos.length > 0 &&
                        this.props.users.currentUser.details.photos.map((image, index) => {
                            return <Card key={index}>
                                <Card.Content>
                                    <div className='ui two buttons'>
                                        <Button basic color='red' onClick={(e,data) => {e.preventDefault(); this.props.removeImage(image);}}>Decline</Button>
                                    </div>
                                </Card.Content>
                                <Image src={image} centered/>
                            </Card>;
                        })
                        }
                    </Card.Group>
                </Segment>
                <Segment secondary textAlign='center'>
                    <Link to={"/users/"} >
                        <Button size='large' content='Cancel'/>
                    </Link>
                    <Button
                        size='large'
                        color='blue'
                        content='Save'
                        onClick={(e,data) => {e.preventDefault(); this.props.saveUser(this.props.users.currentUser)}}
                        />
                </Segment>

                </Form>
            </Segment>
            </Dimmer.Dimmable>
        </Segment.Group>
        )
    }
}
