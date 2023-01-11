import React from 'react'
import ReactDOM from 'react-dom';
import {Provider} from 'react-redux';
import store from '../store/index';
import Routing from '../routes/index';
import 'semantic-ui-css/semantic.min.css';

ReactDOM.render(
    <Provider store={store}>
        <Routing/>
    </Provider>,
    document.getElementById('root')
);
