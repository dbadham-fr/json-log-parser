package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var sampleLog = map[string]interface{}{
	"level":     "debug",
	"timestamp": "2024-06-10T08:47:35.469Z",
	"thread":    "Thread-1",
	"logger":    "Logger-1",
	"message":   "Something bad happened when processing this request - see exception for details",
	"exception": "java.lang.NullPointerException: null\\n\\tat java.base/java.util.Objects.requireNonNull(Unknown Source)\\n\\tat org.forgerock.http.vertx.server.StreamedBodyHandler.<init>(StreamedBodyHandler.java:35)\\n\\tat org.forgerock.http.vertx.server.ChfApplicationWebHandler.handleRequestBody(ChfApplicationWebHandler.java:275)\\n\\tat org.forgerock.http.vertx.server.ChfApplicationWebHandler.handle0(ChfApplicationWebHandler.java:202)\\n\\tat org.forgerock.http.vertx.server.ChfApplicationWebHandler.handle(ChfApplicationWebHandler.java:167)\\n\\tat org.forgerock.http.vertx.server.ChfApplicationWebHandler.handle(ChfApplicationWebHandler.java:93)\\n\\tat org.forgerock.openig.launcher.IdentityGatewayVerticle.adminEndpointRouter(IdentityGatewayVerticle.java:99)\\n\\tat io.vertx.lang.rx.DelegatingHandler.handle(DelegatingHandler.java:20)\\n\\tat io.vertx.core.http.impl.Http1xServerRequestHandler.handle(Http1xServerRequestHandler.java:67)\\n\\tat io.vertx.core.http.impl.Http1xServerRequestHandler.handle(Http1xServerRequestHandler.java:30)\\n\\tat io.vertx.core.impl.ContextImpl.emit(ContextImpl.java:335)\\n\\tat io.vertx.core.impl.DuplicatedContext.emit(DuplicatedContext.java:176)\\n\\tat io.vertx.core.http.impl.Http1xServerConnection.handleMessage(Http1xServerConnection.java:174)\\n\\tat io.vertx.core.net.impl.ConnectionBase.read(ConnectionBase.java:159)\\n\\tat io.vertx.core.net.impl.VertxHandler.channelRead(VertxHandler.java:153)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:442)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:420)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.fireChannelRead(AbstractChannelHandlerContext.java:412)\\n\\tat io.netty.channel.ChannelInboundHandlerAdapter.channelRead(ChannelInboundHandlerAdapter.java:93)\\n\\tat io.netty.handler.codec.http.websocketx.extensions.WebSocketServerExtensionHandler.onHttpRequestChannelRead(WebSocketServerExtensionHandler.java:160)\\n\\tat io.netty.handler.codec.http.websocketx.extensions.WebSocketServerExtensionHandler.channelRead(WebSocketServerExtensionHandler.java:83)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:442)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:420)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.fireChannelRead(AbstractChannelHandlerContext.java:412)\\n\\tat io.vertx.core.http.impl.Http1xUpgradeToH2CHandler.channelRead(Http1xUpgradeToH2CHandler.java:124)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:444)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:420)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.fireChannelRead(AbstractChannelHandlerContext.java:412)\\n\\tat io.netty.handler.codec.ByteToMessageDecoder.fireChannelRead(ByteToMessageDecoder.java:346)\\n\\tat io.netty.handler.codec.ByteToMessageDecoder.channelRead(ByteToMessageDecoder.java:318)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:444)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:420)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.fireChannelRead(AbstractChannelHandlerContext.java:412)\\n\\tat io.vertx.core.http.impl.Http1xOrH2CHandler.end(Http1xOrH2CHandler.java:61)\\n\\tat io.vertx.core.http.impl.Http1xOrH2CHandler.channelRead(Http1xOrH2CHandler.java:38)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:444)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:420)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.fireChannelRead(AbstractChannelHandlerContext.java:412)\\n\\tat io.netty.channel.DefaultChannelPipeline$HeadContext.channelRead(DefaultChannelPipeline.java:1410)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:440)\\n\\tat io.netty.channel.AbstractChannelHandlerContext.invokeChannelRead(AbstractChannelHandlerContext.java:420)\\n\\tat io.netty.channel.DefaultChannelPipeline.fireChannelRead(DefaultChannelPipeline.java:919)\\n\\tat io.netty.channel.nio.AbstractNioByteChannel$NioByteUnsafe.read(AbstractNioByteChannel.java:166)\\n\\tat io.netty.channel.nio.NioEventLoop.processSelectedKey(NioEventLoop.java:788)\\n\\tat io.netty.channel.nio.NioEventLoop.processSelectedKeysOptimized(NioEventLoop.java:724)\\n\\tat io.netty.channel.nio.NioEventLoop.processSelectedKeys(NioEventLoop.java:650)\\n\\tat io.netty.channel.nio.NioEventLoop.run(NioEventLoop.java:562)\\n\\tat io.netty.util.concurrent.SingleThreadEventExecutor$4.run(SingleThreadEventExecutor.java:997)\\n\\tat io.netty.util.internal.ThreadExecutorMap$2.run(ThreadExecutorMap.java:74)\\n\\tat io.netty.util.concurrent.FastThreadLocalRunnable.run(FastThreadLocalRunnable.java:30)\\n\\tat java.base/java.lang.Thread.run(Unknown Source)\\n",
}

var builder = &strings.Builder{}

func BenchmarkLogPrintingWithBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		formatLog(builder, sampleLog)
		builder.Reset()
	}
}

func TestParseSeverity(t *testing.T) {
	var scenarios = map[string]Severity{
		"":            unknown,
		"sdfdfsdfsfd": unknown,
		"unknown":     unknown,
		"trace":       trace,
		"TRACE":       trace,
		"debug":       debug,
		"DEBUG":       debug,
		"info":        info,
		"INFO":        info,
		"warn":        warn,
		"WARN":        warn,
		"error":       err,
		"ERROR":       err,
	}

	for desc := range scenarios {
		assert.Equal(t, scenarios[desc], parseSeverity(desc))
	}

}
