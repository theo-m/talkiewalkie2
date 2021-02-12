/**
 * @fileoverview gRPC-Web generated client stub for tw
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as tw_pb from './tw_pb';


export class ABClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: any; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'binary';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname;
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodInfoGet = new grpcWeb.AbstractClientBase.MethodInfo(
    tw_pb.AddressBook,
    (request: tw_pb.AddressBook) => {
      return request.serializeBinary();
    },
    tw_pb.AddressBook.deserializeBinary
  );

  get(
    request: tw_pb.AddressBook,
    metadata: grpcWeb.Metadata | null): Promise<tw_pb.AddressBook>;

  get(
    request: tw_pb.AddressBook,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: tw_pb.AddressBook) => void): grpcWeb.ClientReadableStream<tw_pb.AddressBook>;

  get(
    request: tw_pb.AddressBook,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: tw_pb.AddressBook) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/tw.AB/Get',
        request,
        metadata || {},
        this.methodInfoGet,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/tw.AB/Get',
    request,
    metadata || {},
    this.methodInfoGet);
  }

}

