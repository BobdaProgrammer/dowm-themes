-- Key mappings configuration
local opts = { noremap = true, silent = true }

-- System clipboard copy and paste
vim.api.nvim_set_keymap('v', '<Leader>y', '"+y', opts)
vim.api.nvim_set_keymap('n', '<Leader>p', '"+p', opts)
vim.api.nvim_set_keymap('v', '<Leader>p', '"+p', opts)

-- Toggle NvimTree
vim.api.nvim_set_keymap('n', '<C-t>', ':NvimTreeToggle<CR>', opts)
