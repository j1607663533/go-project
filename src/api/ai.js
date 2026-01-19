import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export const sendChatMessage = async (message) => {
    try {
        const response = await axios.post(`${API_BASE_URL}/ai/chat`, { message });
        return response.data;
    } catch (error) {
        console.error('Request AI Chat failed:', error);
        throw error;
    }
};
