//
//  CreateQuestionModal.swift
//  Conversa
//
//  Created by zy on 26/07/22.
//

import SwiftUI

struct CreateQuestionModal: View {
    @Binding var isPresented: Bool
    @Binding var title: String
    @Binding var description: String
    
    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            HStack{
                Button(action: {
                    isPresented.toggle()
                    title = ""
                    description = ""
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
            }
            Text("Title")
                .font(.subheadline)
                .bold()
                .foregroundColor(Color(hex: "#383838"))
                .padding(.top, 16)
            TextEditor(text: $title)
                .frame(width: 284, height: 38)
                .padding(.horizontal, 16)
                .foregroundColor(Color(hex: "#383838"))
                .background(Color(hex: "#F4F4F4"))
                .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                .cornerRadius(10)
            
            Text("Description")
                .font(.subheadline)
                .bold()
                .foregroundColor(Color(hex: "#383838"))
                .padding(.top, 10)
            TextEditor(text: $description)
                .frame(width: 284, height: 159)
                .padding(.horizontal, 16)
                .foregroundColor(Color(hex: "#383838"))
                .background(Color(hex: "#F4F4F4"))
                .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                .cornerRadius(10)
            
            Button(action: {
                createPost()
            }, label: {
                Text("Post")
                    .foregroundColor(Color(hex: "#383838"))
                    .frame(
                        minWidth: 0,
                        maxWidth: 314,
                        minHeight: 0,
                        maxHeight: 38
                    )
                    .background(Color(hex: "#FFEBBB"))
                    .cornerRadius(10)
                    .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
            })
        }
        .frame(minHeight: 0, maxHeight: .infinity, alignment: .topLeading)
    }
    
    private func createPost() {
        let jsonData = try? JSONEncoder().encode(CreatePost(title: title, description: description))
        
        // Send request
        guard let url = URL(string: "http://127.0.0.1:8888/post") else {return}
            
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
                title = ""
                description = ""
            } else {
                // Handle unexpected error
                print("unexpected error")
            }
        }
        task.resume()
    }
}

struct CreateQuestionModal_Previews: PreviewProvider {
    static var previews: some View {
        CreateQuestionModal(isPresented: .constant(false), title: .constant(""), description: .constant(""))
    }
}
