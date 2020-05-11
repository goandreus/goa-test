/**
 * @fileoverview gRPC-Web generated client stub for userquery
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.userquery = require('./userquery_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.userquery.UserqueryClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.userquery.UserqueryPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.userquery.SignInRequest,
 *   !proto.userquery.SignInResponse>}
 */
const methodDescriptor_Userquery_SignIn = new grpc.web.MethodDescriptor(
  '/userquery.Userquery/SignIn',
  grpc.web.MethodType.UNARY,
  proto.userquery.SignInRequest,
  proto.userquery.SignInResponse,
  /**
   * @param {!proto.userquery.SignInRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userquery.SignInResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.userquery.SignInRequest,
 *   !proto.userquery.SignInResponse>}
 */
const methodInfo_Userquery_SignIn = new grpc.web.AbstractClientBase.MethodInfo(
  proto.userquery.SignInResponse,
  /**
   * @param {!proto.userquery.SignInRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userquery.SignInResponse.deserializeBinary
);


/**
 * @param {!proto.userquery.SignInRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.userquery.SignInResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.userquery.SignInResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.userquery.UserqueryClient.prototype.signIn =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/userquery.Userquery/SignIn',
      request,
      metadata || {},
      methodDescriptor_Userquery_SignIn,
      callback);
};


/**
 * @param {!proto.userquery.SignInRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.userquery.SignInResponse>}
 *     A native promise that resolves to the response
 */
proto.userquery.UserqueryPromiseClient.prototype.signIn =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/userquery.Userquery/SignIn',
      request,
      metadata || {},
      methodDescriptor_Userquery_SignIn);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.userquery.SignOutRequest,
 *   !proto.userquery.SignOutResponse>}
 */
const methodDescriptor_Userquery_SignOut = new grpc.web.MethodDescriptor(
  '/userquery.Userquery/SignOut',
  grpc.web.MethodType.UNARY,
  proto.userquery.SignOutRequest,
  proto.userquery.SignOutResponse,
  /**
   * @param {!proto.userquery.SignOutRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userquery.SignOutResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.userquery.SignOutRequest,
 *   !proto.userquery.SignOutResponse>}
 */
const methodInfo_Userquery_SignOut = new grpc.web.AbstractClientBase.MethodInfo(
  proto.userquery.SignOutResponse,
  /**
   * @param {!proto.userquery.SignOutRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userquery.SignOutResponse.deserializeBinary
);


/**
 * @param {!proto.userquery.SignOutRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.userquery.SignOutResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.userquery.SignOutResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.userquery.UserqueryClient.prototype.signOut =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/userquery.Userquery/SignOut',
      request,
      metadata || {},
      methodDescriptor_Userquery_SignOut,
      callback);
};


/**
 * @param {!proto.userquery.SignOutRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.userquery.SignOutResponse>}
 *     A native promise that resolves to the response
 */
proto.userquery.UserqueryPromiseClient.prototype.signOut =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/userquery.Userquery/SignOut',
      request,
      metadata || {},
      methodDescriptor_Userquery_SignOut);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.userquery.GetAllSessionsRequest,
 *   !proto.userquery.SessionCollection>}
 */
const methodDescriptor_Userquery_GetAllSessions = new grpc.web.MethodDescriptor(
  '/userquery.Userquery/GetAllSessions',
  grpc.web.MethodType.UNARY,
  proto.userquery.GetAllSessionsRequest,
  proto.userquery.SessionCollection,
  /**
   * @param {!proto.userquery.GetAllSessionsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userquery.SessionCollection.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.userquery.GetAllSessionsRequest,
 *   !proto.userquery.SessionCollection>}
 */
const methodInfo_Userquery_GetAllSessions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.userquery.SessionCollection,
  /**
   * @param {!proto.userquery.GetAllSessionsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userquery.SessionCollection.deserializeBinary
);


/**
 * @param {!proto.userquery.GetAllSessionsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.userquery.SessionCollection)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.userquery.SessionCollection>|undefined}
 *     The XHR Node Readable Stream
 */
proto.userquery.UserqueryClient.prototype.getAllSessions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/userquery.Userquery/GetAllSessions',
      request,
      metadata || {},
      methodDescriptor_Userquery_GetAllSessions,
      callback);
};


/**
 * @param {!proto.userquery.GetAllSessionsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.userquery.SessionCollection>}
 *     A native promise that resolves to the response
 */
proto.userquery.UserqueryPromiseClient.prototype.getAllSessions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/userquery.Userquery/GetAllSessions',
      request,
      metadata || {},
      methodDescriptor_Userquery_GetAllSessions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.userquery.GetIDPURLRequest,
 *   !proto.userquery.GetIDPURLResponse>}
 */
const methodDescriptor_Userquery_GetIDPURL = new grpc.web.MethodDescriptor(
  '/userquery.Userquery/GetIDPURL',
  grpc.web.MethodType.UNARY,
  proto.userquery.GetIDPURLRequest,
  proto.userquery.GetIDPURLResponse,
  /**
   * @param {!proto.userquery.GetIDPURLRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userquery.GetIDPURLResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.userquery.GetIDPURLRequest,
 *   !proto.userquery.GetIDPURLResponse>}
 */
const methodInfo_Userquery_GetIDPURL = new grpc.web.AbstractClientBase.MethodInfo(
  proto.userquery.GetIDPURLResponse,
  /**
   * @param {!proto.userquery.GetIDPURLRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userquery.GetIDPURLResponse.deserializeBinary
);


/**
 * @param {!proto.userquery.GetIDPURLRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.userquery.GetIDPURLResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.userquery.GetIDPURLResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.userquery.UserqueryClient.prototype.getIDPURL =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/userquery.Userquery/GetIDPURL',
      request,
      metadata || {},
      methodDescriptor_Userquery_GetIDPURL,
      callback);
};


/**
 * @param {!proto.userquery.GetIDPURLRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.userquery.GetIDPURLResponse>}
 *     A native promise that resolves to the response
 */
proto.userquery.UserqueryPromiseClient.prototype.getIDPURL =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/userquery.Userquery/GetIDPURL',
      request,
      metadata || {},
      methodDescriptor_Userquery_GetIDPURL);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.userquery.CheckTokenRequest,
 *   !proto.userquery.CheckTokenResponse>}
 */
const methodDescriptor_Userquery_CheckToken = new grpc.web.MethodDescriptor(
  '/userquery.Userquery/CheckToken',
  grpc.web.MethodType.UNARY,
  proto.userquery.CheckTokenRequest,
  proto.userquery.CheckTokenResponse,
  /**
   * @param {!proto.userquery.CheckTokenRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userquery.CheckTokenResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.userquery.CheckTokenRequest,
 *   !proto.userquery.CheckTokenResponse>}
 */
const methodInfo_Userquery_CheckToken = new grpc.web.AbstractClientBase.MethodInfo(
  proto.userquery.CheckTokenResponse,
  /**
   * @param {!proto.userquery.CheckTokenRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userquery.CheckTokenResponse.deserializeBinary
);


/**
 * @param {!proto.userquery.CheckTokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.userquery.CheckTokenResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.userquery.CheckTokenResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.userquery.UserqueryClient.prototype.checkToken =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/userquery.Userquery/CheckToken',
      request,
      metadata || {},
      methodDescriptor_Userquery_CheckToken,
      callback);
};


/**
 * @param {!proto.userquery.CheckTokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.userquery.CheckTokenResponse>}
 *     A native promise that resolves to the response
 */
proto.userquery.UserqueryPromiseClient.prototype.checkToken =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/userquery.Userquery/CheckToken',
      request,
      metadata || {},
      methodDescriptor_Userquery_CheckToken);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.userquery.GetUsersByIDRequest,
 *   !proto.userquery.RegisteredUserCollection>}
 */
const methodDescriptor_Userquery_GetUsersByID = new grpc.web.MethodDescriptor(
  '/userquery.Userquery/GetUsersByID',
  grpc.web.MethodType.UNARY,
  proto.userquery.GetUsersByIDRequest,
  proto.userquery.RegisteredUserCollection,
  /**
   * @param {!proto.userquery.GetUsersByIDRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userquery.RegisteredUserCollection.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.userquery.GetUsersByIDRequest,
 *   !proto.userquery.RegisteredUserCollection>}
 */
const methodInfo_Userquery_GetUsersByID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.userquery.RegisteredUserCollection,
  /**
   * @param {!proto.userquery.GetUsersByIDRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userquery.RegisteredUserCollection.deserializeBinary
);


/**
 * @param {!proto.userquery.GetUsersByIDRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.userquery.RegisteredUserCollection)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.userquery.RegisteredUserCollection>|undefined}
 *     The XHR Node Readable Stream
 */
proto.userquery.UserqueryClient.prototype.getUsersByID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/userquery.Userquery/GetUsersByID',
      request,
      metadata || {},
      methodDescriptor_Userquery_GetUsersByID,
      callback);
};


/**
 * @param {!proto.userquery.GetUsersByIDRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.userquery.RegisteredUserCollection>}
 *     A native promise that resolves to the response
 */
proto.userquery.UserqueryPromiseClient.prototype.getUsersByID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/userquery.Userquery/GetUsersByID',
      request,
      metadata || {},
      methodDescriptor_Userquery_GetUsersByID);
};


module.exports = proto.userquery;

