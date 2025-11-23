<template>
  <div>
    <h1>欢迎使用 DevDesk</h1>
    <p>请选择左侧的功能开始使用。</p>

    <section>
      <h2>功能介绍</h2>

      <h3>代码分享</h3>
      <p>
        填入作者、语言、保存时间，将代码粘贴进去，点击分享，然后复制链接，分享给他人，他人即可看到你的代码。
      </p>

      <h3>待办事项</h3>
      <p>
        点击开始，然后添加你的待办事项，并复制链接，通过链接可再次访问和修改，可以进行团队合作，但请保管好你的链接。
      </p>

      <h3>Markdown 编辑器</h3>
      <p>
        点击开始，左侧为编辑区，右侧可预览，可通过链接协同合作，但请注意需要刷新网页才可看到别人更新的内容，同时该功能支持导出 Markdown、HTML、PDF 格式。
      </p>      <h3>线上代码编辑器（待开发）</h3>
      <p>
        提供在线编写、运行和调试代码的环境，支持多种主流语言和代码高亮、自动保存等功能，当前功能正在开发中。
      </p>

      <h3>在线 HTTP 接口测试（待开发）</h3>
      <p>
        支持在线构造和发送 HTTP 请求，查看响应结果，方便前后端联调和接口调试，当前功能正在开发中。
      </p>
    </section>
  </div>
</template>

<script>
export default {
  name: "MomentWelcome",
};
</script>

<style scoped>
h1 {
  font-size: 24px;
  margin-bottom: 16px;
}

h2 {
  margin-top: 24px;
  font-size: 20px;
}

h3 {
  margin-top: 16px;
  font-size: 16px;
}

p {
  margin: 8px 0;
  line-height: 1.6;
}
</style>
