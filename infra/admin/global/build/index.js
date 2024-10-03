exports.handler = async (event, context) => {
    try {
      const placeholderResponse = {
        message: "Lambda executed successfully!",
        functionName: context.functionName,
        memoryLimitInMB: context.memoryLimitInMB,
        logGroupName: context.logGroupName,
        input: event
      };
  
      return {
        statusCode: 200,
        body: JSON.stringify(placeholderResponse),
      };
    } catch (error) {
      return {
        statusCode: 500,
        body: JSON.stringify({ error: "Something went wrong" }),
      };
    }
  };
  