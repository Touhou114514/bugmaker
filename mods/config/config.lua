return {
  -- 颜色相关
  source_list_line_color = "\27[34m", -- ANSI 34 (深蓝)
  source_list_keyword_color = "\27[0m",
  source_list_string_color = "\27[92m",
  source_list_number_color = "\27[0m",
  source_list_comment_color = "\27[95m",
  source_list_arrow_color = "\27[93m",
  source_list_tab_color = "\27[90m",
  prompt_color = "",
  stacktrace_function_color = "",
  stacktrace_basename_color = "",

  -- 行数与格式
  source_list_line_count = 5,
  tab = "  ",

  -- 路径替换规则
  substitute_path = {
    -- {from = "/old/path", to = "/new/path"}
  },

  -- 别名
  aliases = {
    -- run = {"r", "start"},
    -- quit = {"q", "exit"}
  },

  -- 数组、字符串等限制
  max_array_values = 64,
  max_string_len = 64,
  max_variable_recurse = 1,

  -- 反汇编风格
  disassemble_flavor = "intel",

  -- 调试信息目录
  debug_info_directories = {"/usr/lib/debug/.build-id"},

  -- 其他
  show_location_expr = false,
  trace_show_timestamp = false,
  position = "default"
}