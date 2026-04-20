<style>
table {
  border-spacing: 0;
  border-collapse: collapse;
}

.td_result {
  border: 1px solid #ACBED1;
}

</style>

<template>
  <div style="width: 100%">
    <div>
      <table style="margin-top: 10px">
        <tbody>
        <tr>
          <td>
            <el-input placeholder="视频目录" v-model="videoDir">
              <template #prepend>视频目录</template>
            </el-input>
          </td>
          <td>
            <el-button type="primary" v-on:click="connectWebSocket">合并</el-button>
          </td>
          <td>
            <el-tag size="medium" type="danger">{{ connState }}</el-tag>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
    <div v-if="messages.length" style="border-color:black; margin-top: 10px">
      <table style="border: black;border:1px;width: 1000px" id="dataTable">
        <tbody>
        <tr v-for="message in messages">
          <td class='td_result'>{{ message }}</td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>

export default {
  data() {
    return {
      websocketUrl: 'ws://127.0.0.1/w/file/mergeVideo', // 替换为你的WebSocket服务器URL
      websocket: null,
      messages: [],
      connState: "未连接",
      videoDir: "D:/duanj",
      websocketConnected: false
    };
  },
  computed: {},
  methods: {
    connectWebSocket() {
      if (this.websocketUrl) {
        this.websocket = new WebSocket(this.websocketUrl);
        this.websocket.onopen = () => {
          this.websocketConnected = true;
          console.log('WebSocket连接已建立');
          this.connState = "连接已建立";
          this.websocket.send(this.videoDir);
        };
        this.websocket.onmessage = event => {
          this.connState = "视频合并处理中";
          this.messages.push(event.data);
          console.log(event.data);
        };
        this.websocket.onerror = error => {
          this.connState = "WebSocket错误"
          console.error('WebSocket错误:', error);
        };
        this.websocket.onclose = () => {
          console.log('WebSocket连接已关闭');
          this.connState = "处理完成"
          this.websocketConnected = false;
        };
      } else {
        this.connState = "未提供WebSocket UR"
        console.error('未提供WebSocket URL');
      }
    },
    sendMessage() {
      this.connectWebSocket()
      if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
        this.websocket.send(this.videoDir);
      } else {
        console.error('WebSocket未连接或连接状态不佳');
      }
    },
  },
};
</script>