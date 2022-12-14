# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import real_estate_service_pb2 as real__estate__service__pb2


class RealEstateStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetListings = channel.unary_unary(
                '/api_v1.RealEstate/GetListings',
                request_serializer=real__estate__service__pb2.GetListingsRequest.SerializeToString,
                response_deserializer=real__estate__service__pb2.GetListingsResponse.FromString,
                )


class RealEstateServicer(object):
    """Missing associated documentation comment in .proto file."""

    def GetListings(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_RealEstateServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetListings': grpc.unary_unary_rpc_method_handler(
                    servicer.GetListings,
                    request_deserializer=real__estate__service__pb2.GetListingsRequest.FromString,
                    response_serializer=real__estate__service__pb2.GetListingsResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'api_v1.RealEstate', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class RealEstate(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GetListings(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/api_v1.RealEstate/GetListings',
            real__estate__service__pb2.GetListingsRequest.SerializeToString,
            real__estate__service__pb2.GetListingsResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
