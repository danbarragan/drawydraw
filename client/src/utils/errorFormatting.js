const formatServerError = (error) => {
  try {
    // Try to access the error field that the server populates
    return error.response.data.error;
  } catch (e) {
    // Default to string representation of the error if it's not a valid JSON response
    return String(error);
  }
};

exports.formatServerError = formatServerError;
