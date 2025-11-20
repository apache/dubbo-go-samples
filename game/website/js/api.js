/**
 * API wrapper
 * Using fetch API instead of jQuery
 */

const baseURL = 'http://127.0.0.1:8089/';

/**
 * Unified API request handler
 * @param {string} url - Request URL
 * @param {object} options - Fetch options
 * @returns {Promise} Returns Promise
 */
async function request(url, options = {}) {
    try {
        const response = await fetch(url, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        
        if (data.code !== 0) {
            throw new Error(data.msg || 'Request failed');
        }

        return data;
    } catch (error) {
        throw error;
    }
}

/**
 * User login
 * @param {object} data - Login data {name: string}
 * @returns {Promise} Returns Promise
 */
function login(data) {
    const params = new URLSearchParams({ name: data.name });
    return request(`${baseURL}login?${params.toString()}`, {
        method: 'GET',
    });
}

/**
 * Submit score
 * @param {object} data - Score data {name: string, score: number}
 * @returns {Promise} Returns Promise
 */
function score(data) {
    return request(`${baseURL}score`, {
        method: 'POST',
        body: JSON.stringify({
            name: data.name,
            score: data.score,
        }),
    });
}

/**
 * Get rank
 * @param {object} data - Rank query data {name: string}
 * @returns {Promise} Returns Promise
 */
function rank(data) {
    const params = new URLSearchParams({ name: data.name });
    return request(`${baseURL}rank?${params.toString()}`, {
        method: 'GET',
    });
}
