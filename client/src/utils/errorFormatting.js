const formatServerError = (error) => {
  try {
    // Try to parse the error message the server sent down
    return JSON.parse(error.response.data).error;
  } catch (e) {
    // Default to string representation of the error if it's not a valid JSON response
    return String(error);
  }
};

exports.formatServerError = formatServerError;
