-- Ensure lazy.nvim is installed
local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not vim.loop.fs_stat(lazypath) then
  vim.fn.system({
    "git",
    "clone",
    "--filter=blob:none",
    "--branch=stable", -- latest stable release
    "https://github.com/folke/lazy.nvim.git",
    lazypath,
  })
end

vim.opt.rtp:prepend(lazypath)
require("lazy").setup({
  { 
     'edluffy/hologram.nvim',
    "nvim-tree/nvim-web-devicons",
     'neovim/nvim-lspconfig',          -- LSP configurations
    'hrsh7th/nvim-cmp',              -- Autocompletion plugin
    'hrsh7th/cmp-nvim-lsp',          -- LSP source for nvim-cmp
    'hrsh7th/cmp-buffer',            -- Buffer completions
    'hrsh7th/cmp-path',              -- Path completions
    'hrsh7th/cmp-cmdline',           -- Command line completions
    'L3MON4D3/LuaSnip',	-- Snippets plugin "nvim-tree/nvim-web-devicons",
  },
{
  "mattn/emmet-vim",
  ft = { "html", "css", "javascript", "typescript", "javascriptreact", "typescriptreact" },
  init = function()
    vim.g.user_emmet_leader_key = '<C-y>'
  end,
},
{
    "nvim-treesitter/nvim-treesitter",
    run = ":TSUpdate",
    config = function()
        require'nvim-treesitter.configs'.setup {
            ensure_installed = { "lua", "rasi", "go" }, -- Add 'rasi' to this list
            highlight = {
                enable = true,
                additional_vim_regex_highlighting = false,
            },
        }
    end,
},
  {
    "neovim/nvim-lspconfig",
    config = function()
      -- Removed gopls setup here
    end,
  },
  {
	  "RRethy/vim-hexokinase",
	  build = "make hexokinase", -- compiles the binary
	  config = function()
		  vim.g.Hexokinase_highlighters = { "foreground", "virtual" }
	  end,
  },
  { 
      "ellisonleao/gruvbox.nvim",
      priority = 1000 , config = true
  },
  {
  "luckasRanarison/tailwind-tools.nvim",
  name = "tailwind-tools",
  build = ":UpdateRemotePlugins",
  dependencies = {
    "nvim-treesitter/nvim-treesitter",
    "nvim-telescope/telescope.nvim", -- optional
    "neovim/nvim-lspconfig", -- optional
  },
  opts = {} -- your configuration
  },
  {
    "nvim-tree/nvim-tree.lua",
    config = function()
      require("nvim-tree").setup {}
    end,
  },
  {
    "nvim-lualine/lualine.nvim",
    requires = { "nvim-tree/nvim-web-devicons", opt = true },
    config = function()
      require("lualine").setup {
        options = {
          theme = "catppuccin",
          section_separators = { "", "" },
          component_separators = { "", "" },
        },
      }
    end,
  },
  {
    'romgrk/barbar.nvim',
    requires = { 'nvim-tree/nvim-web-devicons' }
  },
  {
    "hrsh7th/nvim-cmp",
    requires = {
      "hrsh7th/cmp-nvim-lsp",    -- LSP source for nvim-cmp
      "hrsh7th/cmp-buffer",      -- Buffer completions
      "hrsh7th/cmp-path",        -- Path completions
      "hrsh7th/cmp-cmdline",     -- Command-line completions
      "L3MON4D3/LuaSnip",        -- Snippet engine
      "saadparwaiz1/cmp_luasnip" -- Snippet completions
    }
  },
})
require('hologram').setup{
    auto_display = true -- WIP automatic markdown image display, may be prone to breaking
}

require("gruvbox").setup({
  terminal_colors = true, -- add neovim terminal colors
  undercurl = true,
  underline = true,
  bold = true,
  italic = {
    strings = true,
    emphasis = true,
    comments = true,
    operators = false,
    folds = true,
  },
  strikethrough = true,
  invert_selection = false,
  invert_signs = false,
  invert_tabline = false,
  inverse = true, -- invert background for search, diffs, statuslines and errors
  contrast = "", -- can be "hard", "soft" or empty string
  palette_overrides = {},
  overrides = {},
  dim_inactive = false,
  transparent_mode = false,
})

-- Lualine configuration
require("lualine").setup {
  options = {
    theme = "gruvbox",
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
vim.o.background = 'light'
vim.cmd("colorscheme gruvbox")

