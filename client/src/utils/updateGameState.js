import axios from 'axios';

async function UpdateGameState(groupName, playerName, onSuccess, onError) {
  try {
    const response = await axios.get(`/api/get-game-status/${groupName}?playerName=${playerName}`);
    onSuccess(response.data);
  } catch (error) {
    onError(error);
  }
}

export default UpdateGameState;
