import { reduxActions } from '../constants';
import { defaultState } from './initial-states.js';

export function getTheme(state = defaultState, action) {
	switch (action.type){
		case reduxActions.GET_THEME_REQUEST:
			return {...state, fetching: true, fetched: false, error: null}
		case reduxActions.GET_THEME_ERROR:
			return {...state, fetching: false, fetched: false, error: action.error.status}
		case reduxActions.RETURN_THEME:
			return {...state, fetching: false, fetched: true, data: action.theme}
		default:
			return state;
	}
}