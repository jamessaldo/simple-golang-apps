//
//  CreateCommentModal.swift
//  Tangled
//
//  Created by zy on 27/07/22.
//

import SwiftUI

struct CreateCommentModal: View {
    var postId: Int
    @Binding var isPresented: Bool
    @Binding var content: String
    @Binding var post: Post
    
    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            HStack(alignment: .top){
                Button(action: {
                    isPresented.toggle()
                    content = ""
                }) {
                    Text("Cancel")
                        .foregroundColor(Color(hex: "#383838"))
                        .frame(
                            width: 68 , height: 30, alignment: .center
                        )
                        .background(Color(hex: "#FFEBBB"))
                        .cornerRadius(10)
                        .padding(.top, 16)
                        .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                }
                Spacer()
                Button(action: {
                    createComment()
                }) {
                    Text("Answer")
                        .foregroundColor(Color(hex: "#383838"))
                        .frame(
                            width: 68 , height: 30, alignment: .center
                        )
                        .background(Color(hex: "#FFEBBB"))
                        .cornerRadius(10)
                        .padding(.top, 16)
                        .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                }
            }
            Rectangle().overlay{
                HStack(alignment: .top){
                    VStack(alignment: .center){
                        Image(systemName: "person.circle.fill")
                            .resizable()
                            .scaledToFit()
                            .frame(width: 32, height: 32)
                            .foregroundColor(.gray)
                            .padding(.top, 16)
                        DottedLine()
                            .stroke(style: StrokeStyle(lineWidth: 1, dash: [2]))
                            .foregroundColor(Color.black)
                            .frame(maxWidth:1, maxHeight: .infinity)
                    }
                    
                    VStack(alignment: .leading){
                        HStack{
                            VStack(alignment: .leading) {
                                Text(post.creator ?? "Fetch error!")
                                    .font(.caption)
                                    .bold()
                                    .foregroundColor(Color(hex: "#383838"))
                                    .padding(.top, 16)
                                    .multilineTextAlignment(.leading)
                                Text(formatter_toString.string(from:formatter_toDate.date(from: post.created_at ?? "") ?? Date()))
                                    .font(.caption)
                                    .foregroundColor(Color(hex: "#383838"))
                                    .multilineTextAlignment(.leading)
                            }
                        }
                        Text(post.description ?? "Fetch error!")
                            .font(.caption)
                            .foregroundColor(Color(hex: "#484848"))
                            .padding(.top, 8)
                            .multilineTextAlignment(.leading)
                        Text("Replying to @\(post.creator ?? "Error")") { string in
                            string.foregroundColor = Color(hex: "#383838")
                            if let range = string.range(of: "@\(post.creator ?? "Error")") { /// here!
                                string[range].foregroundColor = .blue
                            }
                        }
                        .font(.caption)
                        .bold()
                        .padding(.top, 8)
                    }
                }.frame(maxWidth: .infinity,alignment: .topLeading)
            }
            .foregroundColor(Color.clear)
            .frame(height: 159)
            Rectangle().overlay{
                HStack(alignment: .top){
                    VStack(alignment: .center){
                        Image(systemName: "person.circle.fill")
                            .resizable()
                            .scaledToFit()
                            .frame(width: 32, height: 32)
                            .foregroundColor(.gray)
                    }
                    
                    VStack(alignment: .leading){
                        TextEditor(text: $content)
                            .frame(height: 159)
                            .padding(.horizontal, 16)
                            .foregroundColor(Color(hex: "#383838"))
                            .background(Color(hex: "#F4F4F4"))
                            .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                            .cornerRadius(10)
                    }
                    
                }.frame(maxHeight: .infinity, alignment: .topLeading)
            }
            .foregroundColor(Color.clear)
            .frame(height: 159)
            
            
            
            Spacer()
            
        }
        .frame(
            minWidth: 0,
            maxWidth: 327,
            alignment: .topLeading
        )
        .clipped()
    }
    
    private func createComment() {
        let jsonData = try? JSONEncoder().encode(CreateComment(content: content, post_id: postId))
        
        // Send request
        guard let url = URL(string: "http://127.0.0.1:8888/comment") else {return}
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.httpBody = jsonData
        
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let task = URLSession.shared.dataTask(with: request) { (data, response, error) in
            if let error = error {
                // Handle HTTP request error
                print(error)
            } else if let _ = data {
                // Handle HTTP request response
                isPresented.toggle()
                content = ""
            } else {
                // Handle unexpected error
                print("unexpected error")
            }
        }
        task.resume()
    }
}


struct CreateCommentModal_Previews: PreviewProvider {
    static var previews: some View {
        CreateCommentModal(postId: 0, isPresented: .constant(false), content: .constant(""), post: .constant(Post()))
            .previewInterfaceOrientation(.portrait)
    }
}
