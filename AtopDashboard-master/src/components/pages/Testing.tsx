import React from "react";
import { Form, Input, Button, Select, message } from "antd";
import axios from "axios";
const layout = {
  labelCol: { span: 8 },
  wrapperCol: { span: 16 },
};
const tailLayout = {
  wrapperCol: { offset: 8, span: 16 },
};

export default () => {
  const [form] = Form.useForm();
  const onFinish = (values: object) => {
    async function post(v: object) {
      try {
        const ret = await axios.post("/api/formtest", v);
      } catch (err) {
        message.error(err.message);
      }
    }
    post(values);
  };
  return (
    <div>
      <h3>Testing</h3>
      <Form {...layout} form={form} name="testing-form" onFinish={onFinish}>
        <Form.Item name="note" label="Note" rules={[{ required: true }]}>
          <Input />
        </Form.Item>

        <Form.Item {...tailLayout}>
          <Button type="primary" htmlType="submit">
            Submit
          </Button>
          <Button htmlType="button" onClick={() => console.log("reset")}>
            Reset
          </Button>
          <Button
            type="link"
            htmlType="button"
            onClick={() => console.log("reset")}
          >
            Fill form
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
};
