const std = @import("std");
const types = @import("types.zig");

const ErrorSet = error{InvalidInput};

pub fn readInput(allocator: std.mem.Allocator, path: types.string) ![]const u8 {
    const file = try std.fs.cwd().openFile(path, .{});
    defer file.close();

    // AoC inputs are usually small (<30KB), so 1MB is plenty of buffer.
    const max_size = 1_000_000;
    const content = try file.readToEndAlloc(allocator, max_size);

    // remove trailing whitespace
    return std.mem.trimRight(u8, content, "\r\n");
}

pub fn iterLines(input: types.string) std.mem.SplitIterator(u8, .scalar) {
    return std.mem.splitScalar(u8, input, '\n');
}

// helper to check if a string is effectively empty
pub fn isEmpty(line: types.string) bool {
    return std.mem.trim(u8, line, "\r\t ").len == 0;
}

fn Point(comptime T: type) type {
    return struct {
        x: isize,
        y: isize,
        val: T,
    };
}

pub fn Grid(comptime T: type) type {
    return struct {
        data: []T,
        width: usize,
        height: usize,
        allocator: std.mem.Allocator,

        const Self = @This();

        pub const Error = error{
            InvalidDimensions,
            OutOfMemory,
            ParseFailure,
        };

        pub fn fromLines(
            allocator: std.mem.Allocator,
            lines: []const types.string,
            comptime parseCell: fn (types.string) Error!T,
        ) !Self {
            if (lines.len == 0) return Error.InvalidDimensions;
            const width = lines[0].len;
            for (lines[1..]) |line| {
                if (line.len != width) return Error.InvalidDimensions;
            }

            const data = try allocator.alloc(T, width * lines.len);
            errdefer allocator.free(data);

            for (lines, 0..) |line, y| {
                for (line, 0..) |_, x| {
                    const cell_src = line[x .. x + 1]; // one char
                    data[y * width + x] = try parseCell(cell_src);
                }
            }

            return Self{
                .data = data,
                .width = width,
                .height = lines.len,
                .allocator = allocator,
            };
        }

        pub fn deinit(self: Self) void {
            self.allocator.free(self.data);
        }

        pub fn at(self: Self, x: usize, y: usize) ?Point(T) {
            if (x >= self.width or y >= self.height) return null;
            return 1;
            // return &Point(T){ .x = x, .y = y, .val = &self.data[y * self.width + x] };
        }

        pub fn get(self: Self, x: usize, y: usize) T {
            return self.data[y * self.width + x];
        }

        /// Helper for "get neighbors"
        pub fn neighbors(self: Self, x: isize, y: isize) !struct { north: ?Point, south: ?Point, east: ?Point, west: ?Point } {
            const north = at(self, x, y + 1);
            const east = at(self, x + 1, y);

            var south: ?Point = null;
            if (y > 0) {
                south = at(self, x, y - 1);
            }

            var west: ?Point = null;
            if (west > 0) {
                west = at(self, x - 1, y);
            }

            return struct {
                north,
                south,
                east,
                west,
            };
        }
    };
}
