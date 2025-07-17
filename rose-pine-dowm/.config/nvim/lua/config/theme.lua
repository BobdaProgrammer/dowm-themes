
-- Theme configuration (Catppuccin)
require("catppuccin").setup({
  flavour = "mocha", -- Choose between latte, frappe, macchiato, mocha
  integrations = {
    cmp = true,
    gitsigns = true,
    nvimtree = true,
    treesitter = true,
  },
})
vim.cmd.colorscheme("catppuccin")

-- Customize highlight groups
vim.api.nvim_set_hl(0, 'NormalFloat', { bg = '#313244', fg = '#f2cdcd' })  -- Background and text color for float windows
vim.api.nvim_set_hl(0, 'FloatBorder', { fg = '#94e2d5', bg = '#313244' })  -- Border color for float windows
vim.api.nvim_set_hl(0, 'Pmenu', { bg = '#313244', fg = '#94e2d5' })        -- Background and text color for completion menu
vim.api.nvim_set_hl(0, 'PmenuSel', { bg = '#6272a4', fg = '#ffffff' })    -- Selected item color in completion menu

-- Lualine configuration
require("lualine").setup {
  options = {
    theme = "catppuccin",
    component_separators = '',
    section_separators = { '', '' },
  },
  sections = {
    lualine_a = { { 'mode', separator = { left = '' }, right_padding = 2 } },
    lualine_b = { 'filename', 'branch' },
    lualine_c = { '%=' },
    lualine_y = { 'filetype', 'progress' },
    lualine_z = { { 'location', separator = { right = '' }, left_padding = 2 } },
  },
}
