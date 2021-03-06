import React from 'react';
import { connect } from 'react-redux';
import styled from 'styled-components';
import { PlaylistThunks } from '../actions';
import { SidebarItem, AddPlaylistItem } from './';

const Wrapper = styled.div`
    display: flex;
    overflow-y: auto;
    width: 300px;
    height: 100%;
    background-color: ${props => props.theme.Background2};

    &::-webkit-scrollbar {
        width: 8px;
        background-color: inherit;
        padding-right: 1px;
    }

    &::-webkit-scrollbar-thumb {
        border-radius: 10px;
        -webkit-box-shadow: inset 0 0 6px rgba(0,0,0,.3);
        background-color: ${props => props.theme.Primary1};
    }
`;

const Container = styled.div`
    position: relative;
    width: 100%;
`;

class Sidebar extends React.Component {
    componentWillMount() {
        this.props.dispatch(PlaylistThunks.getAll());
    }

    render() {
        return (
            <Wrapper>
                <Container>
                    {this.props.playlists.map(plst => {
                        return <SidebarItem key={plst.Name} plst={plst} />;
                    })}
                    <AddPlaylistItem />
                </Container>
            </Wrapper>
        );
    }
}

function select(state) {
    return {
        playlists: state.playlists.data
    };
}

export default connect(select)(Sidebar);
